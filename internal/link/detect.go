package link

import (
	"os"

	"github.com/mlaccetti/muster/internal/adapter"
)

// DetectInstalled returns harness IDs resolvable on PATH (ignores MUSTER_ADAPTER_ROOT).
func DetectInstalled() ([]string, error) {
	_ = os.Unsetenv("MUSTER_ADAPTER_ROOT")
	var out []string
	for _, id := range adapter.AllIDs() {
		a, err := adapter.Get(id)
		if err != nil {
			continue
		}
		if _, _, err := a.ResolveCLI(); err == nil {
			out = append(out, id)
		}
	}
	return out, nil
}

func normalizeHarnessList(requested []string) ([]string, error) {
	if len(requested) == 0 {
		return DetectInstalled()
	}
	out := make([]string, 0, len(requested))
	for _, h := range requested {
		id, err := adapter.NormalizeHarnessID(h)
		if err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, nil
}
