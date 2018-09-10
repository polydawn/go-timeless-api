package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/warpfork/go-wish"
)

func TestFilesetFilterSerialization(t *testing.T) {
	atl := atlas.MustBuild(FilesetPackFilter_AtlasEntry)
	shouldRtPack := func(t *testing.T, obj FilesetPackFilter, canon string) {
		t.Helper()
		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := FilesetPackFilter{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	}
	t.Run("pack lossless should roundtrip", func(t *testing.T) {
		shouldRtPack(t,
			FilesetPackFilter_Lossless,
			`"uid=keep,gid=keep,mtime=keep,sticky=keep,setid=keep,dev=keep"`)
	})
	t.Run("pack flatten should roundtrip", func(t *testing.T) {
		shouldRtPack(t,
			FilesetPackFilter_Flatten,
			`"uid=1000,gid=1000,mtime=@1262304000,sticky=keep,setid=keep,dev=keep"`)
	})

	atl = atlas.MustBuild(FilesetUnpackFilter_AtlasEntry)
	shouldRtUnpack := func(t *testing.T, obj FilesetUnpackFilter, canon string) {
		t.Helper()
		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := FilesetUnpackFilter{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	}
	t.Run("unpack lossless should roundtrip", func(t *testing.T) {
		shouldRtUnpack(t,
			FilesetUnpackFilter_Lossless,
			`"uid=keep,gid=keep,mtime=keep,sticky=keep,setid=keep,dev=keep"`)
	})
	t.Run("unpack lowpriv should roundtrip", func(t *testing.T) {
		shouldRtUnpack(t,
			FilesetUnpackFilter_LowPriv,
			`"uid=mine,gid=mine,mtime=keep,sticky=keep,setid=reject,dev=reject"`)
	})
}