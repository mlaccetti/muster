package link

import (
	"errors"
	"fmt"
	"os"

	"github.com/mlaccetti/muster/internal/harnesscli"
	"github.com/mlaccetti/muster/internal/kit"
)

func adapterAllIDs() []string {
	return harnesscli.AllIDs()
}

type Options struct {
	Harnesses []string // empty = auto-detect on PATH (unless All)
	All       bool     // link all known harness IDs without PATH detect
	DryRun    bool
	Unlink    bool
}

type Result struct {
	KitDir    string
	Harnesses []HarnessLink
	Skipped   []string
}

func Run(opts Options) (*Result, error) {
	_ = os.Unsetenv("MUSTER_ADAPTER_ROOT")

	kitDir, err := kit.ActiveDir()
	if err != nil {
		kitDir, err = kit.Install()
		if err != nil {
			return nil, fmt.Errorf("kit not installed: %w (run muster install)", err)
		}
	}

	var ids []string
	var err2 error
	if opts.All {
		ids = append([]string{}, adapterAllIDs()...)
	} else {
		ids, err2 = normalizeHarnessList(opts.Harnesses)
		if err2 != nil {
			return nil, err2
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no harness CLIs found on PATH; install a CLI, pass --harness <id>, or --all")
	}

	res := &Result{KitDir: kitDir}
	autoDetect := !opts.All && len(opts.Harnesses) == 0
	for _, id := range ids {
		h := harnessWiring{id: id}
		if opts.Unlink {
			if err := h.unlinkGlobal(opts.DryRun); err != nil {
				return nil, fmt.Errorf("%s unlink: %w", id, err)
			}
			continue
		}
		var targets []string
		t, err := h.linkGlobal(kitDir, opts.DryRun)
		if err != nil {
			if autoDetect || errors.Is(err, errManualHarnessSetup) {
				res.Skipped = append(res.Skipped, id+": "+err.Error())
				continue
			}
			return nil, fmt.Errorf("%s harness link: %w", id, err)
		}
		targets = append(targets, t...)
		res.Harnesses = append(res.Harnesses, HarnessLink{ID: id, Targets: targets})
	}

	if opts.Unlink {
		_ = removeLinkShims(opts.DryRun)
		if !opts.DryRun {
			_ = os.Remove(manifestPathSafe())
		}
		return res, nil
	}

	if err := installLinkShims(opts.DryRun); err != nil {
		return nil, err
	}
	if !opts.DryRun {
		m := &Manifest{KitDir: kitDir, Harnesses: res.Harnesses}
		if err := SaveManifest(m); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func manifestPathSafe() string {
	p, err := manifestPath()
	if err != nil {
		return ""
	}
	return p
}
