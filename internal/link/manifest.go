package link

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/mlaccetti/muster/internal/paths"
)

const manifestName = "link-manifest.json"

type Manifest struct {
	Version   int           `json:"version"`
	UpdatedAt string        `json:"updated_at"`
	KitDir    string        `json:"kit_dir"`
	Harnesses []HarnessLink `json:"harnesses"`
}

type HarnessLink struct {
	ID      string   `json:"id"`
	Targets []string `json:"targets"`
}

func manifestPath() (string, error) {
	home, err := paths.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, manifestName), nil
}

func LoadManifest() (*Manifest, error) {
	p, err := manifestPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func SaveManifest(m *Manifest) error {
	p, err := manifestPath()
	if err != nil {
		return err
	}
	m.Version = 1
	m.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}
