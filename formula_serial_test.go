package api

import (
	"testing"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestFormulaSerializationFixtures(t *testing.T) {
	t.Run("basic formula", func(t *testing.T) {
		ShouldMarshalPrettyJson(t,
			Formula{
				Inputs: FormulaInputs{
					"/": WareID{"demo", "asdf"},
				},
				Action: FormulaAction{
					Exec: []string{"/bin/hello", "world"},
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
				"outputs": null
			}`,
		)
	})
}
