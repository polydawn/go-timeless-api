package api

import (
	"testing"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestHitchSerializationFixtures(t *testing.T) {
	t.Run("ReleaseItemID serialization", func(t *testing.T) {
		ShouldMarshalJson(t,
			ReleaseItemID{"a", "b", "c"},
			HitchAtlas,
			`"a:b:c"`,
		)
	})
	t.Run("Catalog serialization", func(t *testing.T) {
		t.Run("empty catalog, no releases", func(t *testing.T) {
			ShouldMarshalJson(t,
				Catalog{
					"cname",
					[]ReleaseEntry{},
				},
				HitchAtlas,
				`{"name":"cname","releases":[]}`,
			)
		})
		t.Run("short catalog: one release, no replay", func(t *testing.T) {
			ShouldMarshalPrettyJson(t,
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
				HitchAtlas,
				`
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
				}`,
			)
		})
	})
}
