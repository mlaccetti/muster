package buildversion

import (
	"runtime/debug"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// String returns the CLI release version (e.g. v0.1.0).
func String() string {
	if v := strings.TrimSpace(version); v != "" && v != "dev" {
		return normalizeTag(v)
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if mod := strings.TrimSpace(info.Main.Version); mod != "" && mod != "(devel)" {
			return normalizeTag(mod)
		}
	}
	return "dev"
}

// Commit and Date are set by GoReleaser ldflags for release binaries.
func Commit() string { return commit }

func Date() string { return date }

// Line is the full --version output (release builds include commit and date).
func Line() string {
	v := String()
	if commit != "" && commit != "none" && date != "" && date != "unknown" {
		return v + " (" + commit + ", " + date + ")\n"
	}
	return v + "\n"
}

func normalizeTag(v string) string {
	v = strings.TrimSpace(v)
	if v == "" || v == "dev" {
		return v
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}
