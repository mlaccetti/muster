package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"

	rootkit "github.com/mlaccetti/muster/kit"
)

var (
	mu       sync.Mutex
	compiled = map[string]*jsonschema.Schema{}
)

func ValidateFile(kind string, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return ValidateBytes(kind, data)
}

func ValidateBytes(kind string, data []byte) error {
	sch, err := schemaFor(kind)
	if err != nil {
		return err
	}
	var v any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return sch.Validate(v)
}

func schemaFor(kind string) (*jsonschema.Schema, error) {
	mu.Lock()
	defer mu.Unlock()
	if s, ok := compiled[kind]; ok {
		return s, nil
	}
	name := map[string]string{
		"brief":   "brief.v1.json",
		"report":  "worker-report.v1.json",
		"routing": "routing.v1.json",
		"pools":   "pools.v1.json",
	}[kind]
	if name == "" {
		return nil, fmt.Errorf("unknown schema kind %q", kind)
	}
	data, err := rootkit.Source.ReadFile("schemas/" + name)
	if err != nil {
		return nil, err
	}
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	loc := "https://github.com/mlaccetti/muster/kit/schemas/" + name
	if err := compiler.AddResource(loc, doc); err != nil {
		return nil, err
	}
	sch, err := compiler.Compile(loc)
	if err != nil {
		return nil, err
	}
	compiled[kind] = sch
	return sch, nil
}
