package release

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMaybeNotify_announcesOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", home)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":"v9.9.9"}`))
	}))
	defer srv.Close()

	oldURL := LatestReleaseURL
	oldClient := Client
	LatestReleaseURL = srv.URL
	Client = srv.Client()
	defer func() {
		LatestReleaseURL = oldURL
		Client = oldClient
	}()

	var buf strings.Builder
	MaybeNotify(context.Background(), "v0.1.0", &buf)
	if !strings.Contains(buf.String(), "brew upgrade muster") || !strings.Contains(buf.String(), "v9.9.9") {
		t.Fatalf("msg %q", buf.String())
	}
	buf.Reset()
	MaybeNotify(context.Background(), "v0.1.0", &buf)
	if buf.Len() != 0 {
		t.Fatalf("expected single announce, got %q", buf.String())
	}
}

func TestMaybeNotify_usesCacheWithin24h(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", home)
	_ = saveState(State{
		LastCheck: time.Now(),
		LatestTag: "v9.9.9",
	})

	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		_, _ = w.Write([]byte(`{"tag_name":"v9.9.9"}`))
	}))
	defer srv.Close()
	oldURL := LatestReleaseURL
	oldClient := Client
	LatestReleaseURL = srv.URL
	Client = srv.Client()
	defer func() {
		LatestReleaseURL = oldURL
		Client = oldClient
	}()

	var buf strings.Builder
	MaybeNotify(context.Background(), "v0.1.0", &buf)
	if called {
		t.Fatal("should not call GitHub when cache is fresh")
	}
	if !strings.Contains(buf.String(), "v9.9.9") {
		t.Fatalf("msg %q", buf.String())
	}
}

func TestMaybeNotify_skipsDev(t *testing.T) {
	var buf strings.Builder
	MaybeNotify(context.Background(), "dev", &buf)
	if buf.Len() != 0 {
		t.Fatal("should skip dev builds")
	}
}

func TestCompareTags(t *testing.T) {
	if compareTags("v0.2.0", "v0.1.0") != 1 {
		t.Fatal()
	}
	if compareTags("v0.1.0", "v0.2.0") != -1 {
		t.Fatal()
	}
}

func TestStateRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", home)
	if err := saveState(State{LatestTag: "v1.0.0", LastCheck: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, stateFileName)); err != nil {
		t.Fatal(err)
	}
}
