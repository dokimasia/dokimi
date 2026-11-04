// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package render

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
	"go.dokimi.dev/spec"
	"go.yaml.in/yaml/v3"

	"go.dokimi.dev/dokimi/rules"
)

// The documents of the rendered files. The field of each section has the
// name of the section in the file, and it maps the identifier of each entry
// to the entry.
type (
	claimsDoc struct {
		Claims map[spec.ClaimID]rules.Claim `json:"claims"`
	}
	lawsDoc struct {
		Laws map[spec.LawID]rules.Law `json:"laws"`
	}
	defectsDoc struct {
		Forms map[spec.DefectFormID]rules.DefectForm `json:"forms"`
	}
	gatesDoc struct {
		Gates map[spec.GateID]rules.Gate `json:"gates"`
	}
	unclaimedDoc struct {
		Unclaimed map[string]rules.Unclaimed `json:"unclaimed"`
	}
	referencesDoc struct {
		References map[rules.ReferenceFormID]rules.ReferenceForm `json:"references"`
	}
)

// document is one file of the rules: its name without the extension, and
// the type of its document.
type document struct {
	typ  reflect.Type
	name string
}

// documents are the files of the rules in the order in which Render writes
// them.
var documents = []document{
	{name: "claims", typ: reflect.TypeFor[claimsDoc]()},
	{name: "laws", typ: reflect.TypeFor[lawsDoc]()},
	{name: "defects", typ: reflect.TypeFor[defectsDoc]()},
	{name: "gates", typ: reflect.TypeFor[gatesDoc]()},
	{name: "unclaimed", typ: reflect.TypeFor[unclaimedDoc]()},
	{name: "references", typ: reflect.TypeFor[referencesDoc]()},
}

// Files returns the names of the files of the rules, without their
// extensions, in the order in which Render writes them.
func Files() []string {
	names := make([]string, 0, len(documents))
	for _, d := range documents {
		names = append(names, d.name)
	}
	return names
}

// identifier is the pattern of the identifier of an entry of the spec:
// lower-case ASCII words joined by hyphens.
const identifier = `^[a-z][a-z0-9]*(-[a-z][a-z0-9]*)*$`

// draft is the JSON Schema draft of each schema.
const draft = "https://json-schema.org/draft/2020-12/schema"

// text is the JSON type of a string.
const text = "string"

// schemaOptions are the options of the inference of each schema. A rule
// embeds the runtime entry of the spec without its fields, so the schema of
// the runtime entry has no properties. Each identifier of the spec has the
// pattern of an identifier, and a state and an effect have the values of the
// rules.
func schemaOptions() *jsonschema.ForOptions {
	runtime := &jsonschema.Schema{Type: "object"}
	id := &jsonschema.Schema{Type: text, Pattern: identifier}
	types := map[reflect.Type]*jsonschema.Schema{
		reflect.TypeFor[spec.Claim]():      runtime,
		reflect.TypeFor[spec.Law]():        runtime,
		reflect.TypeFor[spec.DefectForm](): runtime,
		reflect.TypeFor[spec.Gate]():       runtime,
		reflect.TypeFor[rules.State]():     {Type: text, Enum: values(rules.States)},
		reflect.TypeFor[rules.Effect]():    {Type: text, Enum: values(rules.Effects)},
	}
	for _, t := range []reflect.Type{
		reflect.TypeFor[spec.LawID](), reflect.TypeFor[spec.ClaimID](), reflect.TypeFor[spec.GateID](),
		reflect.TypeFor[spec.MachineryID](), reflect.TypeFor[spec.StanceID](),
	} {
		types[t] = id
	}
	return &jsonschema.ForOptions{TypeSchemas: types}
}

// values returns the values of an enumeration of the rules as the values of
// an enum of a schema.
func values[T ~string](enum []T) []any {
	out := make([]any, len(enum))
	for i, v := range enum {
		out[i] = string(v)
	}
	return out
}

// Render renders the YAML file of each name of Files in dir into its JSON
// file and its schema under dir/schemas.
//
// # Errors
//
// Render returns an error that names the file for a file that does not
// parse, for a field that its entry does not have, and for a file that it
// cannot read or write.
func Render(dir string) (err error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("open the directory of the rules: %w", err)
	}
	//dokimi:mutate-skip sbr-delete: no spec makes the close of a root fail, and the close of a directory has no effect that a caller observes
	defer func() { err = errors.Join(err, root.Close()) }()
	return render(root)
}

// render renders the rules in root, as Render states.
func render(root *os.Root) error {
	opts := schemaOptions()
	if err := root.MkdirAll("schemas", 0o750); err != nil {
		return fmt.Errorf("the directory of the schemas: %w", err)
	}
	for _, d := range documents {
		if err := renderFile(root, d); err != nil {
			return err
		}
		schema, err := jsonschema.ForType(d.typ, opts)
		//dokimi:mutate-skip ror-true: ForType returns an error only for a type without a JSON Schema, and each document has one
		if err == nil {
			schema.Schema = draft
			err = writeJSON(root, "schemas/"+d.name+".schema.json", schema)
		}
		if err != nil {
			return fmt.Errorf("the schema of %s: %w", d.name, err)
		}
	}
	return nil
}

// renderFile renders the YAML file of d in root into its JSON file.
func renderFile(root *os.Root, d document) error {
	source := d.name + ".yaml"
	data, err := root.ReadFile(source)
	if err != nil {
		return fmt.Errorf("%s: %w", source, err)
	}
	var file yaml.Node
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("%s: %w", source, err)
	}
	if len(file.Content) == 0 || file.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("%s: the file is no mapping of sections", source)
	}
	top := file.Content[0]
	var raw bytes.Buffer
	raw.WriteByte('{')
	for i := 0; i < len(top.Content); i += 2 {
		key, value := top.Content[i].Value, top.Content[i+1]
		if i > 0 {
			raw.WriteByte(',')
		}
		section, err := entries(key, value)
		if err != nil {
			return fmt.Errorf("%s: section %s: %w", source, key, err)
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
		return fmt.Errorf("%s: %w", source, err)
	}
	return writeJSON(root, d.name+".json", doc)
}

// entries returns the section key, a mapping from the identifier of each
// entry to the entry, as a JSON object.
func entries(key string, section *yaml.Node) ([]byte, error) {
	if section.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("the section %s is no mapping of identifiers to entries", key)
	}
	var out bytes.Buffer
	out.WriteByte('{')
	for i := 0; i < len(section.Content); i += 2 {
		id, entry := section.Content[i].Value, section.Content[i+1]
		if i > 0 {
			out.WriteByte(',')
		}
		data, err := object(entry, id, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		k, _ := json.Marshal(id)
		out.Write(k)
		out.WriteByte(':')
		out.Write(data)
	}
	out.WriteByte('}')
	return out.Bytes(), nil
}

// object returns the entry n as a JSON object, whose members start with the
// member first when first is not nil. Each mapping under the key inputs or
// partners becomes an array of objects with the name of each under the key
// name, and a scalar under the key unless becomes an object with the scalar
// under the key mixin. The error for an entry that is no mapping states id,
// the identifier or the name of the entry.
func object(n *yaml.Node, id string, first []byte) ([]byte, error) {
	if n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("the entry %s is no mapping", id)
	}
	var out bytes.Buffer
	out.WriteByte('{')
	out.Write(first)
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
		if i > 0 || first != nil {
			out.WriteByte(',')
		}
		k, _ := json.Marshal(key)
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
		name := n.Content[i].Value
		v, _ := json.Marshal(name)
		data, err := object(n.Content[i+1], name, append([]byte(`"name":`), v...))
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
