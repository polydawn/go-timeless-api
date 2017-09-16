package api

import (
	"testing"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestFormulaSerializationFixtures(t *testing.T) {
	t.Run("basic formula", func(t *testing.T) {
		ShouldMarshalPrettyJson(t,
			Formula{
				Inputs: UnpackTree{
					"/": UnpackSpec{WareID: WareID{"demo", "asdf"}},
				},
				Action: FormulaAction{
					Exec: []string{"/bin/hello", "world"},
				},
			},
			RepeatrAtlas,
			`
			{
				"inputs": {
					"/": {
						"ware": "demo:asdf",
						"opts": {
							"uid": "",
							"gid": "",
							"mtime": "",
							"sticky": false
						}
					}
				},
				"action": {
					"exec": [
						"/bin/hello",
						"world"
					]
				},
				"outputs": null
			}`,
		)
	})
}
