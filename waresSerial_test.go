package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/warpfork/go-wish"
)

func TestWareIDSerialization(t *testing.T) {
	atl := atlas.MustBuild(WareID_AtlasEntry)
	obj := WareID{"pck", "hsh"}
	canon := `"pck:hsh"`

	t.Run("marshal", func(t *testing.T) {
		bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, string(bs), ShouldEqual, canon)
	})
	t.Run("unmarshal", func(t *testing.T) {
		targ := WareID{}
		err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, targ, ShouldEqual, obj)
	})
}
