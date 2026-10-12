package pools

import (
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

func TestWalkAvoidFamilyAndWrite(t *testing.T) {
	pf, err := Load(testutil.Path(t, "kit/pools.json"))
	if err != nil {
		t.Fatal(err)
	}
	all, err := Walk("implement", pf, nil, true)
	if err != nil || len(all) < 2 {
		t.Fatalf("implement pool: %d err %v", len(all), err)
	}
	avoid, err := Walk("implement", pf, []string{"anthropic"}, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range avoid {
		if c.Family == "anthropic" {
			t.Fatal("anthropic should be avoided")
		}
	}
}
