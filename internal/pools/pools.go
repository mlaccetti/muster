package pools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mlaccetti/muster/internal/schema"
)

type File struct {
	Version   int                    `json:"version"`
	PassOrder []string               `json:"pass_order"`
	Tiers     map[string][]Candidate `json:"tiers"`
}

type Candidate struct {
	Harness  string `json:"harness"`
	Model    string `json:"model"`
	Family   string `json:"family"`
	Write    bool   `json:"write"`
	Boundary string `json:"boundary,omitempty"`
}

func Load(path string) (*File, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := schema.ValidateBytes("pools", data); err != nil {
		return nil, err
	}
	var f File
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

func Walk(tier string, f *File, avoidFamilies []string, requireWrite bool) ([]Candidate, error) {
	list, ok := f.Tiers[tier]
	if !ok {
		return nil, fmt.Errorf("unknown tier %q", tier)
	}
	avoid := map[string]bool{}
	for _, fam := range avoidFamilies {
		avoid[fam] = true
	}
	out := make([]Candidate, 0, len(list))
	for _, c := range list {
		if avoid[c.Family] {
			continue
		}
		if requireWrite && !c.Write {
			continue
		}
		out = append(out, c)
	}
	return out, nil
}
