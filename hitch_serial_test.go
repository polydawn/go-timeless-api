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
		wantStringEqual(t, string(msg), `"a:b:c"`)
	})
	t.Run("Catalog serialization", func(t *testing.T) {
		t.Run("empty catalog, no releases", func(t *testing.T) {
			msg, err := refmt.MarshalAtlased(json.EncodeOptions{},
				Catalog{
					"cname",
					[]ReleaseEntry{},
				},
				HitchAtlas)
			assertNoError(t, err)
			wantStringEqual(t, string(msg), `{"name":"cname","releases":[]}`)
		})
	})
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func wantStringEqual(t *testing.T, a, b string) {
	result, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
	})
	if err != nil {
		t.Fatalf("diffing failed: %s", err)
	}
	t.Helper()
	if result != "" {
		t.Errorf("Match failed: diff:\n%s", result)
	}
}
