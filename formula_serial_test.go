package api

import (
	"testing"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestFormulaSerializationFixtures(t *testing.T) {
	t.Run("basic formula", func(t *testing.T) {
		ShouldMarshalPrettyJson(t,
			Formula{
				Inputs: map[AbsPath]WareID{
					"/": WareID{"demo", "asdf"},
				},
				Action: FormulaAction{
					Exec: []string{"/bin/hello", "world"},
				},
				Outputs: map[AbsPath]OutputSpec{
					"/saveme": {PackType: "tar"},
				},
				FetchUrls: map[AbsPath][]WarehouseAddr{
					"/": []WarehouseAddr{
						"https+ca://ports.polydawn.io/assets/",
						"https+ca://mirror.wahoo.io/timeless/assets/",
					},
				},
				SaveUrls: map[AbsPath]WarehouseAddr{
					"/saveme": "file+ca://./wares/",
				},
			},
			RepeatrAtlas,
			`
			{
				"inputs": {
					"/": "demo:asdf"
				},
				"action": {
					"exec": [
						"/bin/hello",
						"world"
					],
					"cwd": "",
					"env": null,
					"hostname": ""
				},
				"outputs": {
					"/saveme": {
						"packtype": "tar",
						"filters": {
							"uid": "",
							"gid": "",
							"mtime": "",
							"sticky": ""
						}
					}
				},
				"fetchUrls": {
					"/": [
						"https+ca://ports.polydawn.io/assets/",
						"https+ca://mirror.wahoo.io/timeless/assets/"
					]
				},
				"saveUrls": {
					"/saveme": "file+ca://./wares/"
				}
			}`,
		)
	})
}
