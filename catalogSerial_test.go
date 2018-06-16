package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	. "github.com/warpfork/go-wish"
)

func TestCatalogSerialization(t *testing.T) {
	atl := Atlas_Catalog
	t.Run("examplary catalog should roundtrip", func(t *testing.T) {
		obj := ModuleCatalog{
			Name: "froob.org/base",
			Releases: []Release{
				{
					Name:     "v1",
					Items:    map[ItemName]WareID{"linux-amd64": WareID{"tar", "6q7G4hWr"}},
					Metadata: map[string]string{"optional": "foobaring"},
					Hazards:  map[string]string{"facemelting": "true"},
				},
			},
		}
		canon := Dedent(`
			{
				"name": "froob.org/base",
				"releases": [
					{
						"name": "v1",
						"items": {
							"linux-amd64": "tar:6q7G4hWr"
						},
						"metadata": {
							"optional": "foobaring"
						},
						"hazards": {
							"facemelting": "true"
						}
					}
				]
			}
		`)

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{Line: []byte{'\n'}, Indent: []byte{'\t'}}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := ModuleCatalog{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
}
