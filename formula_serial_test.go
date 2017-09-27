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
					"/saveme": {PackFmt: "tar"},
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
					]
				},
				"outputs": {
					"/saveme": {
						"packfmt": "tar",
						"filters": {
							"uid": "",
							"gid": "",
							"mtime": "",
							"sticky": false
						}
					}
				}
			}`,
		)
	})
}
