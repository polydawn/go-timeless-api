package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/warpfork/go-wish"
)

func TestOperationSerialization(t *testing.T) {
	atl := atlas.MustBuild(
		Operation_AtlasEntry,
		SlotRef_AtlasEntry,
		FormulaAction_AtlasEntry,
		FormulaUserinfo_AtlasEntry,
	)
	t.Run("zero operation should roundtrip", func(t *testing.T) {
		obj := Operation{}
		canon := `{"inputs":null,"action":{}}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := Operation{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
		t.Run("unmarshal blank", func(t *testing.T) {
			targ := Operation{}
			canon := `{}`
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
	t.Run("examplary operation should roundtrip", func(t *testing.T) {
		obj := Operation{
			Inputs: map[SlotRef]AbsPath{
				SlotRef{"", "foo"}: "/",
			},
			Action: FormulaAction{
				Exec: []string{"/script"},
			},
			Outputs: map[SlotName]AbsPath{
				"bar": "/bar",
			},
		}
		canon := Dedent(`
			{
				"inputs": {
					"foo": "/"
				},
				"action": {
					"exec": [
						"/script"
					]
				},
				"outputs": {
					"bar": "/bar"
				}
			}
		`)

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{Line: []byte{'\n'}, Indent: []byte{'\t'}}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := Operation{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
}
