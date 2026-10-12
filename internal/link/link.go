package link

import (
	"fmt"
	"os"

	"github.com/mlaccetti/muster/internal/adapter"
	"github.com/mlaccetti/muster/internal/kit"
)

func adapterAllIDs() []string {
	return adapter.AllIDs()
}

type Options struct {
	ProjectDir string
	Harnesses  []string // empty = auto-detect on PATH (unless All)
	All        bool     // link all known harness IDs without PATH detect
	Global     bool
	Project    bool
	DryRun     bool
	Unlink     bool
}

type Result struct {
	KitDir    string
	Project   string
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

	projectDir := opts.ProjectDir
	if projectDir == "" {
		projectDir, err = os.Getwd()
		if err != nil {
			return nil, err
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

	res := &Result{KitDir: kitDir, Project: projectDir}
	autoDetect := !opts.All && len(opts.Harnesses) == 0
	for _, id := range ids {
		h := harnessWiring{id: id}
		if opts.Unlink {
			if opts.Global {
				if err := h.unlinkGlobal(opts.DryRun); err != nil {
					return nil, fmt.Errorf("%s global unlink: %w", id, err)
				}
			}
			if opts.Project {
				if err := h.unlinkProject(projectDir, opts.DryRun); err != nil {
					return nil, fmt.Errorf("%s project unlink: %w", id, err)
				}
			}
			continue
		}
		var targets []string
		if opts.Global {
			t, err := h.linkGlobal(kitDir, opts.DryRun)
			if err != nil {
				if autoDetect {
					res.Skipped = append(res.Skipped, id+": "+err.Error())
					continue
				}
				return nil, fmt.Errorf("%s global link: %w", id, err)
			}
			targets = append(targets, t...)
		}
		if opts.Project {
			t, err := h.linkProject(projectDir, kitDir, opts.DryRun)
			if err != nil {
				return nil, fmt.Errorf("%s project link: %w", id, err)
			}
			targets = append(targets, t...)
		}
		res.Harnesses = append(res.Harnesses, HarnessLink{ID: id, Targets: targets})
	}

	if opts.Unlink {
		if opts.Project && !opts.DryRun {
			_ = unlinkProject(projectDir, false)
		}
		if !opts.DryRun {
			_ = os.Remove(manifestPathSafe())
		}
		return res, nil
	}

	if opts.Project {
		if _, err := linkProject(projectDir, kitDir, opts.DryRun); err != nil {
			return nil, err
		}
	}

	if !opts.DryRun {
		m := &Manifest{KitDir: kitDir, Project: projectDir, Harnesses: res.Harnesses}
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
