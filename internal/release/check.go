package release

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mlaccetti/muster/internal/paths"
)

const (
	githubRepo    = "mlaccetti/muster"
	checkInterval = 24 * time.Hour
	stateFileName = "release-check.json"
	envSkip       = "MUSTER_SKIP_UPDATE_CHECK"
)

type State struct {
	LastCheck    time.Time `json:"last_check"`
	LatestTag    string    `json:"latest_tag"`
	AnnouncedTag string    `json:"announced_tag"`
}

type ghRelease struct {
	TagName string `json:"tag_name"`
}

// Client fetches release metadata (overridable in tests).
var Client = http.DefaultClient

// LatestReleaseURL is the GitHub releases/latest endpoint (overridable in tests).
var LatestReleaseURL = "https://api.github.com/repos/" + githubRepo + "/releases/latest"

// MaybeNotify checks GitHub at most once per day and prints an upgrade hint to w.
func MaybeNotify(ctx context.Context, current string, w io.Writer) {
	if w == nil || os.Getenv(envSkip) == "1" {
		return
	}
	current = strings.TrimSpace(current)
	if current == "" || current == "dev" {
		return
	}
	state, err := loadState()
	if err != nil {
		state = State{}
	}
	now := time.Now()
	latest := state.LatestTag
	if now.Sub(state.LastCheck) >= checkInterval || latest == "" {
		tag, err := fetchLatestTag(ctx)
		if err != nil {
			return
		}
		latest = tag
		state.LastCheck = now
		state.LatestTag = latest
		_ = saveState(state)
	}
	if !isNewer(latest, current) {
		return
	}
	if state.AnnouncedTag == latest {
		return
	}
	_, _ = fmt.Fprintf(w, "muster %s is available (you have %s). Run: brew update && brew upgrade muster\n", latest, normalizeTag(current))
	state.AnnouncedTag = latest
	_ = saveState(state)
}

func fetchLatestTag(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LatestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "muster-cli")
	resp, err := Client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github status %d", resp.StatusCode)
	}
	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	if rel.TagName == "" {
		return "", fmt.Errorf("empty tag")
	}
	return rel.TagName, nil
}

func statePath() (string, error) {
	home, err := paths.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, stateFileName), nil
}

func loadState() (State, error) {
	p, err := statePath()
	if err != nil {
		return State{}, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return State{}, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}
	return s, nil
}

func saveState(s State) error {
	p, err := statePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func normalizeTag(v string) string {
	v = strings.TrimSpace(v)
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}

func isNewer(latest, current string) bool {
	return compareTags(normalizeTag(latest), normalizeTag(current)) > 0
}

// compareTags returns 1 if a>b, -1 if a<b, 0 if equal (simple numeric semver segments).
func compareTags(a, b string) int {
	ap := parseTag(a)
	bp := parseTag(b)
	for i := 0; i < 3; i++ {
		if ap[i] > bp[i] {
			return 1
		}
		if ap[i] < bp[i] {
			return -1
		}
	}
	return 0
}

func parseTag(v string) [3]int {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	parts := strings.Split(v, ".")
	var out [3]int
	for i := 0; i < 3 && i < len(parts); i++ {
		n := 0
		for _, c := range parts[i] {
			if c < '0' || c > '9' {
				break
			}
			n = n*10 + int(c-'0')
		}
		out[i] = n
	}
	return out
}
