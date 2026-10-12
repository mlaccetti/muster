package render

import (
	"os"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

func TestBriefFromFile(t *testing.T) {
	tmpl := testutil.Path(t, "kit/templates/implement.tmpl")
	input := testutil.Path(t, "testdata/briefs/minimal.json")
	out, err := BriefFromFile(tmpl, input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Touch example.txt") {
		t.Fatalf("output: %s", out)
	}
}

func TestBriefFromBytes(t *testing.T) {
	tmplBytes, err := os.ReadFile(testutil.Path(t, "kit/templates/implement.tmpl"))
	if err != nil {
		t.Fatal(err)
	}
	in, err := os.ReadFile(testutil.Path(t, "testdata/briefs/minimal.json"))
	if err != nil {
		t.Fatal(err)
	}
	out, err := BriefFromBytes(string(tmplBytes), in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "verify_command") && !strings.Contains(out, "Verify") {
		t.Fatalf("missing verify section: %s", out)
	}
}
