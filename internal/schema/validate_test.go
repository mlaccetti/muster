package schema

import (
	"path/filepath"
	"testing"
)

func TestValidateBriefFixture(t *testing.T) {
	brief := filepath.Join("..", "..", "testdata", "briefs", "minimal.json")
	if err := ValidateFile("brief", brief); err != nil {
		t.Fatal(err)
	}
}
