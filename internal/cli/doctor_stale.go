package cli

import "path/filepath"

func staleKitLink(manifestKit, activeKit string) bool {
	if manifestKit == "" || activeKit == "" {
		return false
	}
	a, err1 := filepath.EvalSymlinks(manifestKit)
	b, err2 := filepath.EvalSymlinks(activeKit)
	if err1 == nil && err2 == nil {
		return a != b
	}
	return manifestKit != activeKit
}
