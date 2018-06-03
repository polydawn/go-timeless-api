package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	. "github.com/warpfork/go-wish"
)

func TestModuleSerialization(t *testing.T) {
	t.Run("zero module should roundtrip", func(t *testing.T) {
		obj := Module{}
		canon := `{"imports":null,"operations":null}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := Module{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
		t.Run("unmarshal blank", func(t *testing.T) {
			targ := Module{}
			canon := `{}`
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
}
