package routing

import (
	"os"
	"testing"
)

func TestParseEmitRouting(t *testing.T) {
	line := Line{
		Verdict: "ok", Harness: Ptr("claude"), Family: Ptr("anthropic"),
		Model: Ptr(""), Report: "/tmp/r.json", ReasonCode: "none", Reset: 0,
	}
	tmp := filepathJoin(t, "routing-out")
	f, err := os.Create(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if err := Emit(f, line); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	raw, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseLine(string(raw[:len(raw)-1])) // trim newline
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Verdict != "ok" || parsed.Report != line.Report {
		t.Fatalf("parsed %+v", parsed)
	}
}

func filepathJoin(t *testing.T, name string) string {
	return t.TempDir() + "/" + name
}
