package api

import (
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestHitchSerializationFixtures(t *testing.T) {
	t.Run("ReleaseItemID serialization", func(t *testing.T) {
		msg, err := refmt.MarshalAtlased(json.EncodeOptions{},
			ReleaseItemID{"a", "b", "c"},
			HitchAtlas)
		AssertNoError(t, err)
		WantStringEqual(t, string(msg), `"a:b:c"`)
	})
	t.Run("Catalog serialization", func(t *testing.T) {
		t.Run("empty catalog, no releases", func(t *testing.T) {
			msg, err := refmt.MarshalAtlased(json.EncodeOptions{},
				Catalog{
					"cname",
					[]ReleaseEntry{},
				},
				HitchAtlas)
			AssertNoError(t, err)
			WantStringEqual(t, string(msg), `{"name":"cname","releases":[]}`)
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
			AssertNoError(t, err)
			WantStringEqual(t, PrettifyJson(msg), Dedent(`
				{
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
				}
			`))
		})
	})
}
