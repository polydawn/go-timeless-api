package funcs

import (
	"testing"

	. "github.com/warpfork/go-wish"

	. "go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch/mock"
)

func TestPinning(t *testing.T) {
	module := Module{
		Imports: map[SlotName]ImportRef{
			"base": ImportRef_Catalog{"publishing.group/base", "v2018", "bin-linux-amd64"},
			"foo":  ImportRef_Ingest{"git", ".:HEAD"},
			"bar":  ImportRef_Catalog{"publishing.group/bar", "v2.0", "bin-linux-amd64"},
		},
		Operations: map[StepName]StepUnion{
			"stepA": Operation{
				Inputs: map[SlotReference]AbsPath{
					{"", "base"}: "/",
					{"", "foo"}:  "/foo",
					{"", "bar"}:  "/bar",
				},
				Action: OpAction{Exec: []string{"mv", "/foo/thinger", "/out/thinger"}},
				Outputs: map[SlotName]AbsPath{
					"intermediate": "/out",
				},
			},
			"stepB": Module{
				Imports: map[SlotName]ImportRef{
					"base":   ImportRef_Catalog{"publishing.group/base", "v2018", "bin-linux-amd64"},
					"bar":    ImportRef_Catalog{"publishing.group/bar", "v2.2", "bin-linux-amd64"}, // n.b. submodule uses different version of bar than parent; that's allowed.
					"wodget": ImportRef_Parent{"stepA", "intermediate"},
				},
				Operations: map[StepName]StepUnion{
					"op": Operation{
						Inputs: map[SlotReference]AbsPath{
							{"", "base"}:   "/",
							{"", "bar"}:    "/bar",
							{"", "wodget"}: "/src",
						},
						Action: OpAction{Exec: []string{"/bar/tool", "/src", "/out/thinger"}},
						Outputs: map[SlotName]AbsPath{
							"intermediate": "/out",
						},
					},
				},
				Exports: map[ItemName]SlotReference{"barred": {"op", "intermediate"}},
			},
			"stepC": Operation{
				Inputs: map[SlotReference]AbsPath{
					{"", "base"}:        "/",
					{"stepB", "barred"}: "/bar",
				},
				Action: OpAction{Exec: []string{"/bar/thinger"}},
				Outputs: map[SlotName]AbsPath{
					"final": "/bar",
				},
			},
		},
		Exports: map[ItemName]SlotReference{
			"src":             {"", "foo"},
			"bin-linux-amd64": {"stepC", "final"},
		},
	}
	pins, err := ResolvePins(module, mockhitch.Fixture{
		map[ModuleName]ModuleCatalog{
			"publishing.group/base": ModuleCatalog{"publishing.group/base", []Release{
				{Name: "v2018",
					Items: map[ItemName]WareID{
						"bin-linux-amd64": WareID{"tar", "asdflkjgh"},
					}},
			}},
			"publishing.group/bar": ModuleCatalog{"publishing.group/bar", []Release{
				{Name: "v2.0",
					Items: map[ItemName]WareID{
						"bin-linux-amd64": WareID{"tar", "qwer1"},
					}},
				{Name: "v2.2",
					Items: map[ItemName]WareID{
						"bin-linux-amd64": WareID{"tar", "qwer2"},
					}},
			}},
		},
	}.ViewCatalog)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, pins, ShouldEqual, map[SubmoduleSlotReference]WareID{
		// {"", SlotReference{"", "foo"}}:       {}, // FUTURE we'll need the ingest tool to make this something interesting
		{"", SlotReference{"", "base"}}:      {"tar", "asdflkjgh"},
		{"", SlotReference{"", "bar"}}:       {"tar", "qwer1"},
		{"stepB", SlotReference{"", "base"}}: {"tar", "asdflkjgh"},
		{"stepB", SlotReference{"", "bar"}}:  {"tar", "qwer2"},
	})
}