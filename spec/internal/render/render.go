// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"go.dokimi.dev/eidos/lang/naming"
	"go.yaml.in/yaml/v3"

	"go.dokimi.dev/dokimi/spec"
)

// The documents of the rendered files. The field of each section has the
// name of the section in the file.
type (
	lawsDoc struct {
		Laws []spec.Law `json:"laws"`
	}
	claimsDoc struct {
		Claims []spec.Claim `json:"claims"`
	}
	gatesDoc struct {
		Gates []spec.Gate `json:"gates"`
	}
	probesDoc struct {
		Probes []spec.Probe `json:"probes"`
	}
	reasonsDoc struct {
		Reasons []spec.Reason `json:"reasons"`
	}
	defectsDoc struct {
		Forms []spec.DefectForm `json:"forms"`
	}
	machineriesDoc struct {
		Machineries []spec.Machinery `json:"machineries"`
	}
	stancesDoc struct {
		Stances   []spec.Stance    `json:"stances"`
		Unclaimed []spec.Unclaimed `json:"unclaimed"`
	}
)

// document is one file of the spec: its name without the extension, and the
// type of its document.
type document struct {
	typ  reflect.Type
	name string
}

// documents are the files of the spec in the order in which Render writes
// them.
var documents = []document{
	{name: "laws", typ: reflect.TypeFor[lawsDoc]()},
	{name: "claims", typ: reflect.TypeFor[claimsDoc]()},
	{name: "gates", typ: reflect.TypeFor[gatesDoc]()},
	{name: "probes", typ: reflect.TypeFor[probesDoc]()},
	{name: "reasons", typ: reflect.TypeFor[reasonsDoc]()},
	{name: "defects", typ: reflect.TypeFor[defectsDoc]()},
	{name: "machineries", typ: reflect.TypeFor[machineriesDoc]()},
	{name: "stances", typ: reflect.TypeFor[stancesDoc]()},
}

// Files returns the names of the files of the spec, without their
// extensions, in the order in which Render writes them.
func Files() []string {
	names := make([]string, 0, len(documents))
	for _, d := range documents {
		names = append(names, d.name)
	}
	return names
}

// kind is the Go type of the identifiers of one section, which names.gen.go
// declares its constants with, and the word that an error uses for an entry
// of the section. A section without a kind, such as unclaimed, declares no
// constants.
type kind struct {
	goType, word string
}

// kinds contains the kind of each section that declares constants.
var kinds = map[string]kind{
	"laws":        {"LawID", "law"},
	"claims":      {"ClaimID", "claim"},
	"gates":       {"GateID", "gate"},
	"probes":      {"ProbeID", "probe"},
	"reasons":     {"ReasonID", "reason"},
	"forms":       {"DefectFormID", "defect form"},
	"machineries": {"MachineryID", "machinery"},
	"stances":     {"StanceID", "stance"},
}

// identifier is the pattern of the identifier of an entry: lower-case ASCII
// words joined by hyphens.
const identifier = `^[a-z][a-z0-9]*(-[a-z][a-z0-9]*)*$`

// articles are the words that no identifier starts with, so that the
// constant of an identifier reads as a name.
var articles = []string{"a", "an", "the"}

// draft is the JSON Schema draft of each schema.
const draft = "https://json-schema.org/draft/2020-12/schema"

// schemaOptions are the options of the inference of each schema. Each
// identifier type has the pattern of an identifier.
func schemaOptions() *jsonschema.ForOptions {
	id := &jsonschema.Schema{Type: "string", Pattern: identifier}
	types := map[reflect.Type]*jsonschema.Schema{}
	for _, t := range []reflect.Type{
		reflect.TypeFor[spec.LawID](), reflect.TypeFor[spec.ClaimID](), reflect.TypeFor[spec.GateID](),
		reflect.TypeFor[spec.ProbeID](), reflect.TypeFor[spec.ReasonID](), reflect.TypeFor[spec.DefectFormID](),
		reflect.TypeFor[spec.MachineryID](), reflect.TypeFor[spec.StanceID](),
	} {
		types[t] = id
	}
	return &jsonschema.ForOptions{TypeSchemas: types}
}

// Render renders the YAML file of each name of Files in dir into its JSON
// file and its schema under dir/schemas, and writes dir/names.gen.go.
//
// # Errors
//
// Render returns an error that names the file for a file that does not
// parse, for a field that its entry does not have, and for a file that it
// cannot read or write. It returns an error for an identifier that starts
// with an article, and for an identifier that two kinds of entry share,
// because the constants of names.gen.go share one namespace.
func Render(dir string) (err error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("open the directory of the spec: %w", err)
	}
	defer func() { err = errors.Join(err, root.Close()) }()
	return render(root)
}

// render renders the spec in root, as Render states.
func render(root *os.Root) error {
	opts := schemaOptions()
	names := map[string]string{}
	var consts []constant
	if err := root.MkdirAll("schemas", 0o750); err != nil {
		return fmt.Errorf("the directory of the schemas: %w", err)
	}
	for _, d := range documents {
		ids, err := renderFile(root, d)
		if err != nil {
			return err
		}
		schema, err := jsonschema.ForType(d.typ, opts)
		if err == nil {
			schema.Schema = draft
			err = writeJSON(root, "schemas/"+d.name+".schema.json", schema)
		}
		if err != nil {
			return fmt.Errorf("the schema of %s: %w", d.name, err)
		}
		for _, s := range ids {
			k := kinds[s.section]
			for _, id := range s.ids {
				if other, ok := names[id]; ok {
					return fmt.Errorf("the identifier %s names a %s and a %s", id, other, k.word)
				}
				names[id] = k.word
				consts = append(consts, constant{name: naming.Pascal(id), goType: k.goType, id: id})
			}
		}
	}
	return writeNames(root, consts)
}

// sectionIDs are the identifiers of the entries of one section of a file,
// in the order of the file.
type sectionIDs struct {
	section string
	ids     []string
}

// renderFile renders the YAML file of d in root into its JSON file, and
// returns the identifiers of each section that declares constants.
func renderFile(root *os.Root, d document) ([]sectionIDs, error) {
	source := d.name + ".yaml"
	data, err := root.ReadFile(source)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}
	var file yaml.Node
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}
	if len(file.Content) == 0 || file.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: the file is no mapping of sections", source)
	}
	top := file.Content[0]
	var ids []sectionIDs
	var raw bytes.Buffer
	raw.WriteByte('{')
	for i := 0; i < len(top.Content); i += 2 {
		key, value := top.Content[i].Value, top.Content[i+1]
		if i > 0 {
			raw.WriteByte(',')
		}
		section, keys, err := entries(key, value)
		if err != nil {
			return nil, fmt.Errorf("%s: section %s: %w", source, key, err)
		}
		if _, ok := kinds[key]; ok {
			ids = append(ids, sectionIDs{section: key, ids: keys})
		}
		k, _ := json.Marshal(key)
		raw.Write(k)
		raw.WriteByte(':')
		raw.Write(section)
	}
	raw.WriteByte('}')
	doc := reflect.New(d.typ).Interface()
	dec := json.NewDecoder(&raw)
	dec.DisallowUnknownFields()
	if err := dec.Decode(doc); err != nil {
		return nil, fmt.Errorf("%s: %w", source, err)
	}
	if err := writeJSON(root, d.name+".json", doc); err != nil {
		return nil, err
	}
	return ids, nil
}

// entries returns the section key, a mapping from the identifier of each
// entry to the entry, as a JSON array of the entries with the identifier
// under the key id, and the identifiers in the order of the section. The
// inputs and the partners of an entry become arrays with the name of each
// under the key name, and a mixin under the key unless becomes an object
// with the mixin under the key mixin.
func entries(key string, section *yaml.Node) ([]byte, []string, error) {
	if section.Kind != yaml.MappingNode {
		return nil, nil, fmt.Errorf("the section %s is no mapping of identifiers to entries", key)
	}
	var out bytes.Buffer
	var ids []string
	out.WriteByte('[')
	for i := 0; i < len(section.Content); i += 2 {
		id, entry := section.Content[i].Value, section.Content[i+1]
		if article, _, found := strings.Cut(id, "-"); found && slices.Contains(articles, article) {
			return nil, nil, fmt.Errorf("%s: the identifier starts with the article %s", id, article)
		}
		if i > 0 {
			out.WriteByte(',')
		}
		data, err := named(entry, "id", id)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", id, err)
		}
		out.Write(data)
		ids = append(ids, id)
	}
	out.WriteByte(']')
	return out.Bytes(), ids, nil
}

// named returns a mapping as a JSON object that starts with the key field
// and the value name. Each mapping under the key inputs or partners becomes
// an array of objects with the name of each under the key name, and a scalar
// under the key unless becomes an object with the scalar under the key mixin.
func named(n *yaml.Node, field, name string) ([]byte, error) {
	if n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("the entry %s is no mapping", name)
	}
	var out bytes.Buffer
	k, _ := json.Marshal(field)
	v, _ := json.Marshal(name)
	out.WriteByte('{')
	out.Write(k)
	out.WriteByte(':')
	out.Write(v)
	for i := 0; i < len(n.Content); i += 2 {
		key, value := n.Content[i].Value, n.Content[i+1]
		var data []byte
		var err error
		switch {
		case (key == "inputs" || key == "partners") && value.Kind == yaml.MappingNode:
			data, err = namedList(value)
		case key == "unless" && value.Kind == yaml.ScalarNode:
			mixin, _ := json.Marshal(value.Value)
			data = append(append([]byte(`{"mixin":`), mixin...), '}')
		default:
			data, err = plain(value)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		k, _ := json.Marshal(key)
		out.WriteByte(',')
		out.Write(k)
		out.WriteByte(':')
		out.Write(data)
	}
	out.WriteByte('}')
	return out.Bytes(), nil
}

// namedList returns a mapping from names to objects as a JSON array of the
// objects, with the name of each under the key name.
func namedList(n *yaml.Node) ([]byte, error) {
	var out bytes.Buffer
	out.WriteByte('[')
	for i := 0; i < len(n.Content); i += 2 {
		if i > 0 {
			out.WriteByte(',')
		}
		data, err := named(n.Content[i+1], "name", n.Content[i].Value)
		if err != nil {
			return nil, err
		}
		out.Write(data)
	}
	out.WriteByte(']')
	return out.Bytes(), nil
}

// plain returns a YAML value as JSON, with the keys of each mapping in the
// order of the file.
func plain(n *yaml.Node) ([]byte, error) {
	switch n.Kind {
	case yaml.MappingNode:
		var out bytes.Buffer
		out.WriteByte('{')
		for i := 0; i < len(n.Content); i += 2 {
			if i > 0 {
				out.WriteByte(',')
			}
			k, _ := json.Marshal(n.Content[i].Value)
			v, err := plain(n.Content[i+1])
			if err != nil {
				return nil, err
			}
			out.Write(k)
			out.WriteByte(':')
			out.Write(v)
		}
		out.WriteByte('}')
		return out.Bytes(), nil
	case yaml.SequenceNode:
		var out bytes.Buffer
		out.WriteByte('[')
		for i, c := range n.Content {
			if i > 0 {
				out.WriteByte(',')
			}
			v, err := plain(c)
			if err != nil {
				return nil, err
			}
			out.Write(v)
		}
		out.WriteByte(']')
		return out.Bytes(), nil
	default:
		var v any
		err := n.Decode(&v)
		var data []byte
		if err == nil {
			data, err = json.Marshal(v)
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", n.Line, err)
		}
		return data, nil
	}
}

// writeJSON writes v to the file name in root as indented JSON with a final
// newline.
func writeJSON(root *os.Root, name string, v any) error {
	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("create %s: %w", name, err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return errors.Join(enc.Encode(v), f.Close())
}

// constant is one constant of names.gen.go.
type constant struct {
	name, goType, id string
}

// header is the start of names.gen.go, up to the first constant.
const header = `// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Code generated by spec/internal/render. DO NOT EDIT.

package spec

// The identifiers of the entries of the spec.
const (
`

// writeNames writes the constants to names.gen.go in root as a formatted Go
// file of the package spec.
func writeNames(root *os.Root, consts []constant) error {
	var b strings.Builder
	b.WriteString(header)
	for _, c := range consts {
		fmt.Fprintf(&b, "\t%s %s = %q\n", c.name, c.goType, c.id)
	}
	b.WriteString(")\n")
	src, err := format.Source([]byte(b.String()))
	if err == nil {
		err = root.WriteFile("names.gen.go", src, 0o600)
	}
	if err != nil {
		return fmt.Errorf("names.gen.go: %w", err)
	}
	return nil
}
