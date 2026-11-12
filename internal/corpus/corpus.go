// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package corpus judges the corpus of a language by the expected outcomes of
// its fixtures. Each fixture has one file of expected outcomes under
// corpus/expect, in the key form of the spec. The overlay of a language
// declares each difference of the language from the expected outcomes, and
// [Judge] fails each difference that the overlay does not declare.
package corpus

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"go.dokimi.dev/spec"
	"go.yaml.in/yaml/v3"

	"go.dokimi.dev/dokimi/internal/report"
)

// ErrFile is the error of a file of expected outcomes or of an overlay that
// does not decode, and of a file of expected outcomes whose fixture is not
// its path.
var ErrFile = errors.New("corpus: the file is invalid")

// Expectation is the expected outcome of one fixture, as its file under
// corpus/expect states it.
type Expectation struct {
	// Fixture is the path of the fixture in the corpus of each language, such
	// as iface/mixin/ttl.
	Fixture string `yaml:"fixture"`

	// Declarations are the expected outcomes of each declaration under the
	// directive suite, by the key form of its name.
	Declarations map[string]Declaration `yaml:"declarations"`

	// Subjects are the subjects of the fixture.
	Subjects []Subject `yaml:"subjects"`

	// Controls are the controls that the corpus applies the suites of the
	// fixture to, such as iface/mixin/ttl/deadline.Store.
	Controls []string `yaml:"controls"`

	// Fakes are the verified fakes that the generator fake writes for the
	// declarations of the fixture, such as
	// iface/contract/persister/persistertest.ContractFake.
	Fakes []string `yaml:"fakes,omitempty"`

	// Level is the conformance level of the fixture, from 1 to 4.
	Level int `yaml:"level"`
}

// fileKeys are the keys of a file of expected outcomes in the order in which
// a person reads them.
var fileKeys = []string{"fixture", "level", "subjects", "controls", "fakes", "declarations"}

// Subject is a subject of a fixture: the directory of its file relative to
// the corpus, a dot and the name of its type, with the keys of the checks
// that the subject excuses with a reason.
type Subject struct {
	Type    string   `yaml:"type"`
	Excused []string `yaml:"excused,omitempty"`
}

// Declaration is the expected outcome of one declaration under the
// directive suite. Each key is in the key form of the spec, such as
// mixed.read/zero-on-error.
type Declaration struct {
	// Checks are the derived checks by their keys.
	Checks map[string]Check `yaml:"checks"`

	// Authored are the keys of the checks of the directive check.
	Authored []string `yaml:"authored"`

	// Ceilings are the keys of the ceilings of the budgets.
	Ceilings []string `yaml:"ceilings"`

	// Refusals are the gates of the refused claims by the keys of their
	// checks.
	Refusals map[string]spec.GateID `yaml:"refusals"`

	// Unclaimed are the classifications without a check.
	Unclaimed []Unclaimed `yaml:"unclaimed"`
}

// Check is the expected outcome of one derived check: its claim, the key of
// the callable of each partner by the name of the partner, its law, and the
// reason code of its probe.
type Check struct {
	Claim    spec.ClaimID      `yaml:"claim"`
	Partners map[string]string `yaml:"partners,omitempty"`
	Law      spec.LawID        `yaml:"law"`
	Reason   spec.ReasonID     `yaml:"reason"`
}

// Unclaimed is an expected classification of a callable without a check.
type Unclaimed struct {
	// Classification is shape:<name>, mixin:<name> or
	// contract:<name>/<role>.
	Classification string `yaml:"classification"`

	// Callable is the key of the callable, such as mixed.put.
	Callable string `yaml:"callable"`

	Stance spec.StanceID `yaml:"stance"`
}

// Overlay is the overlay of one language, corpus/<language>/overlay.yaml:
// the differences that the language declares, and the highest conformance
// level that the gate judges.
type Overlay struct {
	Language string  `yaml:"language"`
	Diverge  []Entry `yaml:"diverge"`
	Level    int     `yaml:"level"`
}

// Entry declares one difference of a language from the expected outcomes of
// a fixture. Expect is the expected kind of the key, and Actual is the kind
// of the language with its gate or its stance: check, authored, ceiling,
// refused, unclaimed or absent. An entry with the stance blocked and no key
// declares that the language cannot state the fixture. An entry with the
// actual kind absent and the stance blocked declares a key that the runtime
// of the language cannot check. Why is the reason why the language differs.
type Entry struct {
	Fixture string        `yaml:"fixture"`
	Key     string        `yaml:"key,omitempty"`
	Expect  string        `yaml:"expect,omitempty"`
	Actual  string        `yaml:"actual,omitempty"`
	Gate    spec.GateID   `yaml:"gate,omitempty"`
	Stance  spec.StanceID `yaml:"stance,omitempty"`
	Why     string        `yaml:"why"`
}

// Phase is a phase of the gate.
type Phase string

// The phases of the gate, in the order in which [Judge] reports them.
const (
	PhaseGeneration    Phase = "generation"
	PhaseFalsification Phase = "falsification"
	PhaseSubject       Phase = "subject"
	PhaseControls      Phase = "controls"
	PhaseFakes         Phase = "fakes"
	PhaseDeclaration   Phase = "declaration"
)

// Outcome is the outcome of one fixture in a run: its lines of the lock and
// the suite records of its declarations.
type Outcome struct {
	Fixture string
	Lock    []string
	Records []report.Record
}

// Finding is one failure of a phase of a fixture, with the expected and the
// actual outcome. Key is the key of a check, of a callable with its
// classification, of a declaration, or the type of a subject or of a
// control, and it is empty for the fixture itself.
type Finding struct {
	Fixture string
	Phase   Phase
	Key     string
	Want    string
	Got     string
}

// ReadExpectations returns the expectation of each file of fsys whose name
// ends in .yaml, in the lexical order of the files.
//
// # Errors
//
// ReadExpectations returns an error that wraps [ErrFile] and names the file
// for a file that does not decode, for a field that Expectation does not
// have, and for a fixture that is not the path of its file without the
// extension.
func ReadExpectations(fsys fs.FS) ([]Expectation, error) {
	var out []Expectation
	err := fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("corpus: %w", err)
		}
		if d.IsDir() || path.Ext(name) != ".yaml" {
			return nil
		}
		var e Expectation
		if err := decode(fsys, name, &e); err != nil {
			return err
		}
		if want := strings.TrimSuffix(name, ".yaml"); e.Fixture != want {
			return fmt.Errorf("%w: %s has the fixture %q, and its path is %s", ErrFile, name, e.Fixture, want)
		}
		out = append(out, e)
		return nil
	})
	return out, err
}

// ReadOverlay returns the overlay of the file name of fsys.
//
// # Errors
//
// ReadOverlay returns an error that wraps [ErrFile] and names the file for a
// file that does not decode and for a field that Overlay does not have.
func ReadOverlay(fsys fs.FS, name string) (Overlay, error) {
	var o Overlay
	err := decode(fsys, name, &o)
	return o, err
}

// decode decodes the YAML file name of fsys into v, and refuses a field that
// v does not have.
func decode(fsys fs.FS, name string, v any) error {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return fmt.Errorf("corpus: %w", err)
	}
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	if err := d.Decode(v); err != nil {
		return fmt.Errorf("%w: %s: %w", ErrFile, name, err)
	}
	return nil
}

// WriteExpectation writes e into its file under dir, the path of its fixture
// with the extension .yaml, as [ReadExpectations] reads it. The new file
// keeps the comment lines that open the file that it replaces, and writes the
// fixture, its level, its subjects, its controls, its fakes and its
// declarations in that order. A fixture without a fake has no key fakes.
func WriteExpectation(dir string, e *Expectation) error {
	name := filepath.Join(dir, filepath.FromSlash(e.Fixture)+".yaml")
	old, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("corpus: %w", err)
	}
	// YAML encodes every value of an Expectation, and a bytes.Buffer takes
	// every write.
	var doc yaml.Node
	_ = doc.Encode(e)
	pairs := slices.Collect(slices.Chunk(doc.Content, 2))
	slices.SortFunc(pairs, func(a, b []*yaml.Node) int {
		return slices.Index(fileKeys, a[0].Value) - slices.Index(fileKeys, b[0].Value)
	})
	doc.Content = slices.Concat(pairs...)
	var b bytes.Buffer
	b.WriteString(head(old))
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	_ = enc.Encode(&doc)
	// The file exists, so it keeps its mode.
	if err := os.WriteFile(name, b.Bytes(), 0o600); err != nil {
		return fmt.Errorf("corpus: %w", err)
	}
	return nil
}

// head returns the comment lines that open data, with the empty lines among
// and after them.
func head(data []byte) string {
	var b strings.Builder
	for line := range strings.Lines(string(data)) {
		if !strings.HasPrefix(line, "#") && strings.TrimSpace(line) != "" {
			break
		}
		b.WriteString(line)
	}
	return b.String()
}

// Fixtures returns the fixtures of the corpus fsys: each directory in a group
// of the fixtures of expected, as its path, sorted. The group of a fixture is
// the directory of its path, such as iface/mixin for iface/mixin/ttl. A group
// that the corpus lacks has no fixture.
//
// # Errors
//
// Fixtures returns an error that wraps [syscall.ENOTDIR] and names the group
// for a group that is not a directory.
func Fixtures(fsys fs.FS, expected []Expectation) ([]string, error) {
	groups := map[string]bool{}
	for i := range expected {
		groups[path.Dir(expected[i].Fixture)] = true
	}
	var out []string
	for _, group := range slices.Sorted(maps.Keys(groups)) {
		info, err := fs.Stat(fsys, group)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		// On Windows, fs.ReadDir reads a file as a directory without entries.
		if err == nil && !info.IsDir() {
			err = &fs.PathError{Op: "readdir", Path: group, Err: syscall.ENOTDIR}
		}
		var entries []fs.DirEntry
		if err == nil {
			entries, err = fs.ReadDir(fsys, group)
		}
		if err != nil {
			return nil, fmt.Errorf("corpus: %w", err)
		}
		for _, d := range entries {
			if d.IsDir() {
				out = append(out, group+"/"+d.Name())
			}
		}
	}
	return out, nil
}

// Outcomes groups the lines of a lock and the suite records of a run by
// fixture: one outcome for each of fixtures, and one for each package outside
// them. module is the module path of the corpus, which Outcomes removes from
// the package of each line and of each record and from the subject of each
// record. A package belongs to the fixture whose path is the package or
// contains it. Outcomes skips the comment lines and the empty lines of the
// lock.
func Outcomes(module string, fixtures, lock []string, records []report.Record) []Outcome {
	out := make([]Outcome, 0, len(fixtures))
	index := make(map[string]int, len(fixtures))
	for _, f := range fixtures {
		index[f] = len(out)
		out = append(out, Outcome{Fixture: f})
	}
	of := func(pkg string) int {
		for _, f := range fixtures {
			if strings.HasPrefix(pkg, f+"/") {
				return index[f]
			}
		}
		// A package that is a fixture, or that an earlier line of its own
		// added, has its index.
		i, ok := index[pkg]
		if !ok {
			i, index[pkg] = len(out), len(out)
			out = append(out, Outcome{Fixture: pkg})
		}
		return i
	}
	relative := func(p string) string { return strings.TrimPrefix(p, module+"/") }
	for _, text := range lock {
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		pkg, rest, _ := strings.Cut(text, "\t")
		pkg = relative(pkg)
		i := of(pkg)
		out[i].Lock = append(out[i].Lock, pkg+"\t"+rest)
	}
	for i := range records {
		r := records[i]
		r.Suite.Package, r.Subject = relative(r.Suite.Package), relative(r.Subject)
		at := of(r.Suite.Package)
		out[at].Records = append(out[at].Records, r)
	}
	return out
}
