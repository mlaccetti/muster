package link

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
	tomlast "github.com/pelletier/go-toml/v2/unstable"
)

const (
	tomlMarkerBegin        = "# MUSTER_LINK_BEGIN"
	tomlMarkerEnd          = "# MUSTER_LINK_END"
	tomlAgentsBegin        = "# MUSTER_LINK_AGENTS_BEGIN"
	tomlAgentsEnd          = "# MUSTER_LINK_AGENTS_END"
	codexConfigStateSuffix = ".muster-link-state.json"
)

var managedCodexValues = map[string]map[string]string{
	"agents": {"default_subagent_model": "gpt-6-luna", "default_subagent_reasoning_effort": "medium"},
}

type codexSavedValue struct {
	Present bool
	Raw     string
}
type codexConfigState struct {
	Version int
	Values  map[string]map[string]codexSavedValue
}
type codexSection struct {
	name       string
	start, end int
	kind       tomlast.Kind
}
type codexAssignment struct {
	section, key string
	value        tomlast.Range
	line         int
}
type codexDocument struct {
	source      string
	sections    []codexSection
	assignments []codexAssignment
}
type sourceEdit struct {
	start, end  int
	replacement string
}

func parseCodexDocument(source string) (*codexDocument, error) {
	var semantic map[string]any
	if err := toml.Unmarshal([]byte(source), &semantic); err != nil {
		return nil, fmt.Errorf("parse Codex config TOML: %w", err)
	}
	doc := &codexDocument{source: source}
	var parser tomlast.Parser
	parser.Reset([]byte(source))
	active := ""
	for parser.NextExpression() {
		n := parser.Expression()
		switch n.Kind {
		case tomlast.Table, tomlast.ArrayTable:
			parts, err := astKeyPath(n.Key())
			if err != nil {
				return nil, err
			}
			if len(parts) == 0 {
				return nil, errors.New("empty TOML table path")
			}
			start := lineStart(source, int(firstKeyOffset(n.Key())))
			if len(doc.sections) > 0 {
				doc.sections[len(doc.sections)-1].end = start
			}
			active = strings.Join(parts, ".")
			doc.sections = append(doc.sections, codexSection{name: active, start: start, end: len(source), kind: n.Kind})
		case tomlast.KeyValue:
			parts, err := astKeyPath(n.Key())
			if err != nil {
				return nil, err
			}
			if len(parts) == 1 {
				doc.assignments = append(doc.assignments, codexAssignment{section: active, key: parts[0], value: n.Value().Raw, line: lineStart(source, int(firstKeyOffset(n.Key())))})
			}
		}
	}
	if err := parser.Error(); err != nil {
		return nil, fmt.Errorf("parse Codex config TOML syntax: %w", err)
	}
	return doc, nil
}

func firstKeyOffset(it tomlast.Iterator) uint32 {
	if it.Next() {
		return it.Node().Raw.Offset
	}
	return 0
}

func astKeyPath(it tomlast.Iterator) ([]string, error) {
	var out []string
	for it.Next() {
		n := it.Node()
		raw := string(n.Data)
		if len(raw) >= 2 && raw[0] == '"' {
			decoded, err := strconv.Unquote(raw)
			if err != nil {
				return nil, fmt.Errorf("decode TOML key %q: %w", raw, err)
			}
			raw = decoded
		} else if len(raw) >= 2 && raw[0] == '\'' {
			raw = raw[1 : len(raw)-1]
		}
		out = append(out, raw)
	}
	return out, nil
}

func lineStart(source string, offset int) int {
	if offset > len(source) {
		offset = len(source)
	}
	if offset < 0 {
		offset = 0
	}
	if i := strings.LastIndexByte(source[:offset], '\n'); i >= 0 {
		return i + 1
	}
	return 0
}
func lineEnd(source string, start int) int {
	if start >= len(source) {
		return len(source)
	}
	if i := strings.IndexByte(source[start:], '\n'); i >= 0 {
		return start + i + 1
	}
	return len(source)
}
func sectionByName(doc *codexDocument, name string) *codexSection {
	for i := range doc.sections {
		if doc.sections[i].name == name && doc.sections[i].kind == tomlast.Table {
			return &doc.sections[i]
		}
	}
	return nil
}
func assignmentsIn(doc *codexDocument, section, key string) []codexAssignment {
	var out []codexAssignment
	for _, a := range doc.assignments {
		if a.section == section && a.key == key {
			out = append(out, a)
		}
	}
	return out
}
func savedConfigState(doc *codexDocument) codexConfigState {
	state := codexConfigState{Version: 1, Values: map[string]map[string]codexSavedValue{}}
	for section, fields := range managedCodexValues {
		state.Values[section] = map[string]codexSavedValue{}
		for key := range fields {
			matches := assignmentsIn(doc, section, key)
			if len(matches) == 0 {
				state.Values[section][key] = codexSavedValue{}
				continue
			}
			r := matches[len(matches)-1].value
			state.Values[section][key] = codexSavedValue{Present: true, Raw: doc.source[r.Offset : r.Offset+r.Length]}
		}
	}
	return state
}

func patchCodexAgentConfig(path string, dryRun bool) error {
	writePath, mode, err := resolveCodexConfigTarget(path)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	source := string(data)
	if _, err := parseCodexDocument(source); err != nil {
		return err
	}
	source = cleanLegacyCodexConfig(source)
	doc, err := parseCodexDocument(source)
	if err != nil {
		return err
	}
	statePath := path + codexConfigStateSuffix
	state := codexConfigState{}
	if saved, readErr := os.ReadFile(statePath); readErr == nil {
		if err := json.Unmarshal(saved, &state); err != nil {
			return fmt.Errorf("read Codex link state: %w", err)
		}
		if state.Version != 1 {
			return fmt.Errorf("unsupported Codex link state version %d", state.Version)
		}
	} else if os.IsNotExist(readErr) {
		state = savedConfigState(doc)
	} else {
		return readErr
	}
	updates := map[string]map[string]string{
		"agents": {"default_subagent_model": strconv.Quote("gpt-6-luna"), "default_subagent_reasoning_effort": strconv.Quote("medium")},
	}
	patched, err := editCodexValues(doc, updates)
	if err != nil {
		return err
	}
	if _, err := parseCodexDocument(patched); err != nil {
		return fmt.Errorf("Codex routing changes would produce invalid TOML: %w", err)
	}
	if dryRun {
		return nil
	}
	stateBytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := writeFileAtomic(statePath, append(stateBytes, '\n'), 0o600); err != nil {
		return err
	}
	return writeFileAtomic(writePath, []byte(patched), mode)
}

func unlinkCodexAgentConfig(path string, dryRun bool) error {
	writePath, mode, err := resolveCodexConfigTarget(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	source := cleanLegacyCodexConfig(string(data))
	doc, err := parseCodexDocument(source)
	if err != nil {
		return err
	}
	statePath := path + codexConfigStateSuffix
	state := codexConfigState{Version: 1, Values: map[string]map[string]codexSavedValue{}}
	hasState := false
	if saved, readErr := os.ReadFile(statePath); readErr == nil {
		hasState = true
		if err := json.Unmarshal(saved, &state); err != nil {
			return fmt.Errorf("read Codex link state: %w", err)
		}
		if state.Version != 1 {
			return fmt.Errorf("unsupported Codex link state version %d", state.Version)
		}
	} else if !os.IsNotExist(readErr) {
		return readErr
	}
	if !hasState {
		cleaned := cleanLegacyCodexConfig(string(data))
		if cleaned == string(data) || dryRun {
			return nil
		}
		if _, err := parseCodexDocument(cleaned); err != nil {
			return fmt.Errorf("legacy Codex config cleanup would produce invalid TOML: %w", err)
		}
		return writeFileAtomic(writePath, []byte(cleaned), mode)
	}
	updates, removals := map[string]map[string]string{}, map[string][]string{}
	for section, fields := range managedCodexValues {
		for key := range fields {
			saved := state.Values[section][key]
			if saved.Present {
				if updates[section] == nil {
					updates[section] = map[string]string{}
				}
				updates[section][key] = saved.Raw
			} else {
				removals[section] = append(removals[section], key)
			}
		}
	}
	patched, err := editCodexValues(doc, updates)
	if err != nil {
		return err
	}
	doc, err = parseCodexDocument(patched)
	if err != nil {
		return fmt.Errorf("restore Codex config: %w", err)
	}
	patched, err = removeCodexValues(doc, removals)
	if err != nil {
		return err
	}
	if _, err := parseCodexDocument(patched); err != nil {
		return fmt.Errorf("Codex unlink would produce invalid TOML: %w", err)
	}
	if dryRun {
		return nil
	}
	if err := writeFileAtomic(writePath, []byte(patched), mode); err != nil {
		return err
	}
	if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func resolveCodexConfigTarget(path string) (string, os.FileMode, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return path, 0o600, nil
	}
	if err != nil {
		return "", 0, err
	}
	target := path
	if info.Mode()&os.ModeSymlink != 0 {
		target, err = filepath.EvalSymlinks(path)
		if err != nil {
			return "", 0, fmt.Errorf("resolve Codex config symlink %s: %w", path, err)
		}
		info, err = os.Stat(target)
		if err != nil {
			return "", 0, err
		}
	}
	if !info.Mode().IsRegular() {
		return "", 0, fmt.Errorf("Codex config %s is not a regular file", path)
	}
	return target, info.Mode().Perm(), nil
}

func cleanLegacyCodexConfig(source string) string {
	source = removeTomlCommentBlock(source, tomlMarkerBegin, tomlMarkerEnd)
	return removeTomlCommentBlock(source, tomlAgentsBegin, tomlAgentsEnd)
}

func removeTomlCommentBlock(source, begin, end string) string {
	var parser tomlast.Parser
	parser.KeepComments = true
	parser.Reset([]byte(source))
	var pending int = -1
	var edits []sourceEdit
	for parser.NextExpression() {
		n := parser.Expression()
		if n.Kind != tomlast.Comment {
			continue
		}
		line := lineStart(source, int(n.Raw.Offset))
		switch string(n.Data) {
		case begin:
			if pending < 0 {
				pending = line
			}
		case end:
			if pending >= 0 {
				edits = append(edits, sourceEdit{start: pending, end: lineEnd(source, int(n.Raw.Offset))})
				pending = -1
			}
		}
	}
	if parser.Error() != nil || len(edits) == 0 {
		return source
	}
	return applySourceEdits(source, edits)
}

func editCodexValues(doc *codexDocument, updates map[string]map[string]string) (string, error) {
	var edits []sourceEdit
	for section, fields := range updates {
		table := sectionByName(doc, section)
		if table == nil {
			var body strings.Builder
			body.WriteString("\n[" + section + "]\n")
			for _, key := range sortedKeys(fields) {
				value := fields[key]
				body.WriteString(key + " = " + value + "\n")
			}
			edits = append(edits, sourceEdit{start: len(doc.source), end: len(doc.source), replacement: body.String()})
			continue
		}
		for _, key := range sortedKeys(fields) {
			value := fields[key]
			matches := assignmentsIn(doc, section, key)
			if len(matches) > 1 {
				return "", fmt.Errorf("Codex config has duplicate %s.%s assignments", section, key)
			}
			if len(matches) == 1 {
				r := matches[0].value
				edits = append(edits, sourceEdit{start: int(r.Offset), end: int(r.Offset + r.Length), replacement: value})
			} else {
				insert := table.end
				if insert < 0 || insert > len(doc.source) {
					return "", errors.New("invalid TOML table source range")
				}
				prefix := ""
				if insert > 0 && doc.source[insert-1] != '\n' {
					prefix = "\n"
				}
				edits = append(edits, sourceEdit{start: insert, end: insert, replacement: prefix + key + " = " + value + "\n"})
			}
		}
	}
	return applySourceEdits(doc.source, edits), nil
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func removeCodexValues(doc *codexDocument, removals map[string][]string) (string, error) {
	var edits []sourceEdit
	for section, keys := range removals {
		for _, key := range keys {
			matches := assignmentsIn(doc, section, key)
			if len(matches) > 1 {
				return "", fmt.Errorf("Codex config has duplicate %s.%s assignments", section, key)
			}
			if len(matches) == 0 {
				continue
			}
			a := matches[0]
			edits = append(edits, sourceEdit{start: a.line, end: lineEnd(doc.source, int(a.value.Offset+a.value.Length))})
		}
	}
	return applySourceEdits(doc.source, edits), nil
}

func applySourceEdits(source string, edits []sourceEdit) string {
	for i := 0; i < len(edits); i++ {
		for j := i + 1; j < len(edits); j++ {
			if edits[j].start > edits[i].start {
				edits[i], edits[j] = edits[j], edits[i]
			}
		}
	}
	for _, e := range edits {
		source = source[:e.start] + e.replacement + source[e.end:]
	}
	return source
}
