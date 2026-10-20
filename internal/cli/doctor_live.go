package cli

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/mlaccetti/muster/internal/diagnostic"
	"github.com/mlaccetti/muster/internal/link"
	"github.com/spf13/cobra"
)

const (
	maxLiveOutput  = 4 << 20
	maxLiveTimeout = 10 * time.Minute
)

func runLiveDoctor(cmd *cobra.Command, codexBin string, timeout time.Duration) (diagnostic.LiveResult, error) {
	var empty diagnostic.LiveResult
	if timeout <= 0 || timeout > maxLiveTimeout {
		return empty, fmt.Errorf("timeout must be positive and at most %s", maxLiveTimeout)
	}
	bin, err := exec.LookPath(codexBin)
	if err != nil {
		return empty, errors.New("codex executable unavailable")
	}
	var nonceBytes [16]byte
	if _, err := rand.Read(nonceBytes[:]); err != nil {
		return empty, errors.New("could not generate check nonce")
	}
	nonce := hex.EncodeToString(nonceBytes[:])
	dir, err := os.MkdirTemp("", "muster-doctor-live-")
	if err != nil {
		return empty, errors.New("could not create temporary working directory")
	}
	defer func() { _ = os.RemoveAll(dir) }()

	ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
	defer cancel()
	prompt := "Perform a read-only delegation diagnostic. Start the child with fresh minimal context; do not fork parent history. You MUST use the native spawn_agent tool to delegate this bounded task to a child agent: reply with this exact nonce and no other text: " + nonce + ". The child must not use tools, access files or the network, or make changes. Prefer the muster-worker agent; if unavailable, use a built-in worker explicitly configured as gpt-6-luna with medium reasoning. Wait for the child to finish. Return only a short status."
	args := []string{"exec", "--json", "--sandbox", "read-only", "--skip-git-repo-check", "-C", dir, "-c", `approval_policy="never"`, prompt}
	var stdout, stderr boundedOutput
	process := exec.CommandContext(ctx, bin, args...)
	process.WaitDelay = 500 * time.Millisecond
	if err := configureLiveProcess(process); err != nil {
		return empty, err
	}
	stop := func() {
		cancel()
		if process.Process != nil && process.Cancel != nil {
			_ = process.Cancel()
		}
	}
	stdout.cancel, stderr.cancel = stop, stop
	process.Stdout, process.Stderr = &stdout, &stderr
	err = process.Run()
	if stdout.exceeded || stderr.exceeded {
		return empty, errors.New("a Codex output stream exceeded the 4 MiB limit")
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return empty, errors.New("codex check timed out")
	}
	if err != nil {
		return empty, errors.New("codex check process failed")
	}
	events := stdout.Bytes()
	rootID, err := diagnostic.RootThreadID(events)
	if err != nil {
		return empty, err
	}
	home := os.Getenv("CODEX_HOME")
	if home == "" {
		home, err = os.UserHomeDir()
		if err != nil {
			return empty, errors.New("could not locate Codex session store")
		}
		home = filepath.Join(home, ".codex")
	}
	rootPath, err := findRollout(home, rootID)
	if err != nil {
		return empty, errors.New("codex root session evidence unavailable")
	}
	result, err := diagnostic.ValidateLiveEvidence(events, rootPath, home, rootID, nonce)
	return result, err
}

type boundedOutput struct {
	mu       sync.Mutex
	data     []byte
	cancel   context.CancelFunc
	exceeded bool
}

func (b *boundedOutput) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	remaining := maxLiveOutput - len(b.data)
	if len(p) > remaining {
		if remaining > 0 {
			b.data = append(b.data, p[:remaining]...)
		}
		b.exceeded = true
		b.cancel()
		return len(p), nil
	}
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *boundedOutput) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]byte(nil), b.data...)
}

func validateCodexLiveRouting() error {
	m, err := link.LoadManifest()
	if err != nil {
		return errors.New("codex link manifest is unavailable")
	}
	home := os.Getenv("CODEX_HOME")
	if home == "" {
		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return errors.New("codex home is unavailable")
		}
		home = filepath.Join(userHome, ".codex")
	}
	wantAgent := filepath.Join(home, "agents", "muster-worker.toml")
	for _, harness := range m.Harnesses {
		if harness.ID != "codex" {
			continue
		}
		targeted := false
		for _, target := range harness.Targets {
			if target == wantAgent {
				targeted = true
				break
			}
		}
		if !targeted {
			return errors.New("codex muster-worker target is absent from link manifest")
		}
		onlyCodex := *m
		onlyCodex.Harnesses = []link.HarnessLink{harness}
		if len(link.VerifyReports(&onlyCodex)) != 0 {
			return errors.New("codex muster-worker wiring failed verification")
		}
		return nil
	}
	return errors.New("codex is absent from link manifest")
}

func findRollout(home, id string) (string, error) {
	if !diagnostic.IsThreadID(id) {
		return "", errors.New("invalid root session id")
	}
	matches, err := filepath.Glob(filepath.Join(home, "sessions", "*", "*", "*", "rollout-*-"+id+".jsonl"))
	if err != nil || len(matches) != 1 {
		return "", errors.New("root rollout not found")
	}
	return matches[0], nil
}
