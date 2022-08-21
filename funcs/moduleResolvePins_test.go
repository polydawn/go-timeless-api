package funcs

import (
	"context"
	"testing"

	. "github.com/warpfork/go-wish"

	. "github.com/polydawn/go-timeless-api"
	mockhitch "github.com/polydawn/go-timeless-api/hitch/mock"
)

func TestPinning(t *testing.T) {
	module := Module{
		Imports: map[SlotName]ImportRef{
			"base": ImportRef_Catalog{"publishing.group/base", "v2018", "bin-linux-amd64"},
			"foo":  ImportRef_Ingest{"git", ".:HEAD"},
			"bar":  ImportRef_Catalog{"publishing.group/bar", "v2.0", "bin-linux-amd64"},
		},
		Steps: map[StepName]StepUnion{
			"stepA": Operation{
				Inputs: map[AbsPath]SlotRef{
					"/":    {"", "base"},
					"/foo": {"", "foo"},
					"/bar": {"", "bar"},
				},
				Action: FormulaAction{Exec: []string{"mv", "/foo/thinger", "/out/thinger"}},
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
				Steps: map[StepName]StepUnion{
					"op": Operation{
						Inputs: map[AbsPath]SlotRef{
							"/":    {"", "base"},
							"/bar": {"", "bar"},
							"/src": {"", "wodget"},
						},
						Action: FormulaAction{Exec: []string{"/bar/tool", "/src", "/out/thinger"}},
						Outputs: map[SlotName]AbsPath{
							"intermediate": "/out",
						},
					},
				},
				Exports: map[ItemName]SlotRef{"barred": {"op", "intermediate"}},
			},
			"stepC": Operation{
				Inputs: map[AbsPath]SlotRef{
					"/":    {"", "base"},
					"/bar": {"stepB", "barred"},
				},
				Action: FormulaAction{Exec: []string{"/bar/thinger"}},
				Outputs: map[SlotName]AbsPath{
					"final": "/bar",
				},
			},
		},
		Exports: map[ItemName]SlotRef{
			"src":             {"", "foo"},
			"bin-linux-amd64": {"stepC", "final"},
		},
	}
	pins, _, err := ResolvePins(
		module,
		mockhitch.Fixture{
			map[ModuleName]Lineage{
				"publishing.group/base": Lineage{"publishing.group/base", []Release{
					{Name: "v2018",
						Items: map[ItemName]WareID{
							"bin-linux-amd64": WareID{"tar", "asdflkjgh"},
						}},
				}},
				"publishing.group/bar": Lineage{"publishing.group/bar", []Release{
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
		}.ViewLineage,
		func(_ context.Context, _ ModuleName) (*WareSourcing, error) {
			return &WareSourcing{}, nil
		},
		func(_ context.Context, ingestRef ImportRef_Ingest) (*WareID, *WareSourcing, error) {
			return &WareID{"git", "f00f"}, &WareSourcing{}, nil
		},
	)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, pins, ShouldEqual, Pins{
		{"", SlotRef{"", "foo"}}:       {"git", "f00f"},
		{"", SlotRef{"", "base"}}:      {"tar", "asdflkjgh"},
		{"", SlotRef{"", "bar"}}:       {"tar", "qwer1"},
		{"stepB", SlotRef{"", "base"}}: {"tar", "asdflkjgh"},
		{"stepB", SlotRef{"", "bar"}}:  {"tar", "qwer2"},
	})
}
