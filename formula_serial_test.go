package api

import (
	"testing"

	. "go.polydawn.net/go-timeless-api/testutil"
)

func TestFormulaSerializationFixtures(t *testing.T) {
	t.Run("basic formula and context", func(t *testing.T) {
		ShouldMarshalPrettyJson(t,
			FormulaUnion{
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
				},
				&FormulaContext{
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
			},
			RepeatrAtlas,
			`
			{
				"formula": {
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
						"userinfo": null,
						"cradle": "",
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
					}
				},
				"context": {
					"fetchUrls": {
						"/": [
							"https+ca://ports.polydawn.io/assets/",
							"https+ca://mirror.wahoo.io/timeless/assets/"
						]
					},
					"saveUrls": {
						"/saveme": "file+ca://./wares/"
					}
				}
			}`,
		)
	})
}

// Test Userinfo serialization.
// Some of these sanity assertions are moderately nontrivial beacuse we
// consider it important to correctly round-trip the unset/default values,
// which for integers we implement as some pointer jiggery.
func TestFormulaUserinfoSerialization(t *testing.T) {
	t.Run("userinfo should serialize", func(t *testing.T) {
		ShouldMarshalPrettyJson(t,
			FormulaUserinfo{},
			RepeatrAtlas,
			`
			{
				"uid": null,
				"gid": null,
				"username": "",
				"homedir": ""
			}`,
		)
	})
	t.Run("an empty object should deserialize as the zero userinfo", func(t *testing.T) {
		ShouldUnmarshalJson(t,
			`{}`,
			RepeatrAtlas,
			FormulaUserinfo{},
		)
	})
	t.Run("a serial object with UID should deserialize into userinfo correctly", func(t *testing.T) {
		i0 := 0
		ShouldUnmarshalJson(t,
			`{"uid":0}`,
			RepeatrAtlas,
			FormulaUserinfo{Uid: &i0},
		)
	})
	t.Run("userinfo with uid should serialize", func(t *testing.T) {
		i0 := 0
		ShouldMarshalPrettyJson(t,
			FormulaUserinfo{Uid: &i0},
			RepeatrAtlas,
			`
			{
				"uid": 0,
				"gid": null,
				"username": "",
				"homedir": ""
			}`,
		)
	})
}
