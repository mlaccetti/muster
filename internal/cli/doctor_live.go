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
	"time"

	"github.com/mlaccetti/muster/internal/diagnostic"
	"github.com/mlaccetti/muster/internal/link"
	"github.com/spf13/cobra"
)

const maxLiveOutput = 4 << 20
const maxLiveTimeout = 10 * time.Minute

func runLiveDoctor(cmd *cobra.Command, codexBin string, timeout time.Duration) (diagnostic.LiveResult, error) {
	var empty diagnostic.LiveResult
	if timeout <= 0 || timeout > maxLiveTimeout {
		return empty, fmt.Errorf("timeout must be positive and at most %s", maxLiveTimeout)
	}
	bin, err := exec.LookPath(codexBin)
	if err != nil {
		return empty, errors.New("Codex executable unavailable")
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
	defer os.RemoveAll(dir)
	captureDir, err := os.MkdirTemp("", "muster-doctor-live-output-")
	if err != nil {
		return empty, errors.New("could not create output capture directory")
	}
	defer os.RemoveAll(captureDir)

	ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
	defer cancel()
	prompt := "Perform a read-only delegation diagnostic. You MUST use the native spawn_agent tool to delegate this bounded task to a child agent: reply with this exact nonce and no other text: " + nonce + ". The child must not use tools, access files or the network, or make changes. Prefer the muster-worker agent; if unavailable, use a built-in worker explicitly configured as gpt-6-luna with medium reasoning. Wait for the child to finish. Return only a short status."
	args := []string{"exec", "--json", "--sandbox", "read-only", "--skip-git-repo-check", "-C", dir, "-c", `approval_policy="never"`, prompt}
	stdoutFile, err := os.Create(filepath.Join(captureDir, "events.jsonl"))
	if err != nil {
		return empty, errors.New("could not create event capture")
	}
	defer stdoutFile.Close()
	stderrFile, err := os.Create(filepath.Join(captureDir, "diagnostic.log"))
	if err != nil {
		return empty, errors.New("could not create diagnostic capture")
	}
	defer stderrFile.Close()
	process := exec.CommandContext(ctx, bin, args...)
	if err := configureLiveProcess(process); err != nil {
		return empty, err
	}
	process.Stdout, process.Stderr = stdoutFile, stderrFile
	err = process.Run()
	if ctx.Err() != nil {
		return empty, errors.New("Codex check timed out")
	}
	if err != nil {
		return empty, errors.New("Codex check process failed")
	}
	if oversized(stdoutFile) || oversized(stderrFile) {
		return empty, errors.New("Codex check output exceeded safety limit")
	}
	if _, err := stdoutFile.Seek(0, 0); err != nil {
		return empty, errors.New("could not read Codex event capture")
	}
	events, err := os.ReadFile(stdoutFile.Name())
	if err != nil {
		return empty, errors.New("could not read Codex event capture")
	}
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
		return empty, errors.New("Codex root session evidence unavailable")
	}
	result, err := diagnostic.ValidateLiveEvidence(events, rootPath, home, rootID, nonce)
	return result, err
}

func validateCodexLiveRouting() error {
	m, err := link.LoadManifest()
	if err != nil {
		return errors.New("Codex link manifest is unavailable")
	}
	home := os.Getenv("CODEX_HOME")
	if home == "" {
		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return errors.New("Codex home is unavailable")
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
			return errors.New("Codex muster-worker target is absent from link manifest")
		}
		onlyCodex := *m
		onlyCodex.Harnesses = []link.HarnessLink{harness}
		if len(link.VerifyReports(&onlyCodex)) != 0 {
			return errors.New("Codex muster-worker wiring failed verification")
		}
		return nil
	}
	return errors.New("Codex is absent from link manifest")
}

func oversized(f *os.File) bool {
	info, err := f.Stat()
	return err != nil || info.Size() > maxLiveOutput
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
