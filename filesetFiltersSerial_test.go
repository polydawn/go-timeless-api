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
	t.Run("pack filters", func(t *testing.T) {
		shouldRtPack := func(t *testing.T, atl atlas.Atlas, obj FilesetPackFilter, canon string) {
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
		t.Run("string", func(t *testing.T) {
			atl := atlas.MustBuild(FilesetPackFilter_AsString_AtlasEntry)
			shouldRtPack := func(t *testing.T, obj FilesetPackFilter, canon string) {
				shouldRtPack(t, atl, obj, canon)
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
			t.Run("pack partial specifications should roundtrip", func(t *testing.T) {
				shouldRtPack(t,
					FilesetPackFilter{true, ff_unspecified, 1000, ff_unspecified, ff_ignore, ff_unspecified, ff_unspecified},
					`"gid=1000,sticky=ignore"`)
			})
		})
		t.Run("obj", func(t *testing.T) {
			atl := atlas.MustBuild(FilesetPackFilter_AtlasEntry)
			shouldRtPack := func(t *testing.T, obj FilesetPackFilter, canon string) {
				shouldRtPack(t, atl, obj, canon)
			}
			t.Run("pack lossless should roundtrip", func(t *testing.T) {
				shouldRtPack(t,
					FilesetPackFilter_Lossless,
					`{"dev":"keep","gid":"keep","mtime":"keep","setid":"keep","sticky":"keep","uid":"keep"}`)
			})
			t.Run("pack flatten should roundtrip", func(t *testing.T) {
				shouldRtPack(t,
					FilesetPackFilter_Flatten,
					`{"dev":"keep","gid":"1000","mtime":"@1262304000","setid":"keep","sticky":"keep","uid":"1000"}`)
			})
			t.Run("pack partial specifications should roundtrip", func(t *testing.T) {
				shouldRtPack(t,
					FilesetPackFilter{true, ff_unspecified, 1000, ff_unspecified, ff_ignore, ff_unspecified, ff_unspecified},
					`{"gid":"1000","sticky":"ignore"}`)
			})
		})
	})

	t.Run("unpack filters", func(t *testing.T) {
		shouldRtUnpack := func(t *testing.T, atl atlas.Atlas, obj FilesetUnpackFilter, canon string) {
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
		t.Run("string", func(t *testing.T) {
			atl := atlas.MustBuild(FilesetUnpackFilter_AsString_AtlasEntry)
			shouldRtUnpack := func(t *testing.T, obj FilesetUnpackFilter, canon string) {
				shouldRtUnpack(t, atl, obj, canon)
			}
			t.Run("unpack lossless should roundtrip", func(t *testing.T) {
				shouldRtUnpack(t,
					FilesetUnpackFilter_Lossless,
					`"uid=follow,gid=follow,mtime=follow,sticky=follow,setid=follow,dev=follow"`)
			})
			t.Run("unpack lowpriv should roundtrip", func(t *testing.T) {
				shouldRtUnpack(t,
					FilesetUnpackFilter_LowPriv,
					`"uid=mine,gid=mine,mtime=follow,sticky=follow,setid=reject,dev=reject"`)
			})
			t.Run("unpack partial specifications should roundtrip", func(t *testing.T) {
				shouldRtUnpack(t,
					FilesetUnpackFilter{true, ff_unspecified, 1000, ff_unspecified, ff_ignore, ff_unspecified, ff_unspecified},
					`"gid=1000,sticky=ignore"`)
			})
		})
		t.Run("obj", func(t *testing.T) {
			atl := atlas.MustBuild(FilesetUnpackFilter_AtlasEntry)
			shouldRtUnpack := func(t *testing.T, obj FilesetUnpackFilter, canon string) {
				shouldRtUnpack(t, atl, obj, canon)
			}
			t.Run("unpack lossless should roundtrip", func(t *testing.T) {
				shouldRtUnpack(t,
					FilesetUnpackFilter_Lossless,
					`{"dev":"follow","gid":"follow","mtime":"follow","setid":"follow","sticky":"follow","uid":"follow"}`)
			})
			t.Run("unpack lowpriv should roundtrip", func(t *testing.T) {
				shouldRtUnpack(t,
					FilesetUnpackFilter_LowPriv,
					`{"dev":"reject","gid":"mine","mtime":"follow","setid":"reject","sticky":"follow","uid":"mine"}`)
			})
			t.Run("unpack partial specifications should roundtrip", func(t *testing.T) {
				shouldRtUnpack(t,
					FilesetUnpackFilter{true, ff_unspecified, 1000, ff_unspecified, ff_ignore, ff_unspecified, ff_unspecified},
					`{"gid":"1000","sticky":"ignore"}`)
			})
		})
	})

}
