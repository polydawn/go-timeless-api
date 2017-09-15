package api

import (
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
)

func TestHitchSerializationFixtures(t *testing.T) {
	t.Run("ReleaseItemID serialization", func(t *testing.T) {
		msg, err := refmt.MarshalAtlased(json.EncodeOptions{},
			ReleaseItemID{"a", "b", "c"},
			HitchAtlas)
		assertNoError(t, err)
		assertStringEqual(t, string(msg), `"a:b:c"`)
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func assertStringEqual(t *testing.T, a, b string) {
	t.Helper()
	result, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
	})
	if result != "" {
		t.Errorf("Match failed: diff:\n%s", result)
	}
}
