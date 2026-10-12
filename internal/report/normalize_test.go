package report

import (
	"strings"
	"testing"
)

func TestNormalizeJSONPassthrough(t *testing.T) {
	in := []byte(`{"verdict":"DONE","changed":[],"commands":[],"evidence":[],"assumptions":[],"not_done":[],"observations":[]}`)
	out, err := NormalizeBytes(in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(out), "{") {
		t.Fatal("expected json")
	}
}

func TestNormalizeMarkdownBlock(t *testing.T) {
	in := []byte("Summary\n\n```json\n{\"verdict\":\"DONE\",\"changed\":[],\"commands\":[],\"evidence\":[],\"assumptions\":[],\"not_done\":[],\"observations\":[]}\n```\n")
	out, err := NormalizeBytes(in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "DONE") {
		t.Fatalf("out %s", out)
	}
}
