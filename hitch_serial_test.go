package api

import (
	"bytes"
	stdjson "encoding/json"
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
		t.Run("short catalog: one release, no replay", func(t *testing.T) {
			msg, err := refmt.MarshalAtlased(json.EncodeOptions{},
				Catalog{
					"cname",
					[]ReleaseEntry{
						{"1.0",
							map[ItemName]WareID{
								"item-a": {"war", "asdf"},
								"item-b": {"war", "qwer"},
							},
							map[string]string{
								"comment": "yes",
							},
							nil,
							nil,
						},
					},
				},
				HitchAtlas)
			assertNoError(t, err)
			wantStringEqual(t, jsonPretty(msg),
				`{
	"name": "cname",
	"releases": [
		{
			"name": "1.0",
			"items": {
				"item-a": "war:asdf",
				"item-b": "war:qwer"
			},
			"metadata": {
				"comment": "yes"
			},
			"hazards": null,
			"replay": null
		}
	]
}`)
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

func jsonPretty(s []byte) string {
	var out bytes.Buffer
	stdjson.Indent(&out, s, "", "\t")
	return out.String()
}
