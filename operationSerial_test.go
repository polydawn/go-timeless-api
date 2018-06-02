package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/warpfork/go-wish"
)

func TestOpActionSerialization(t *testing.T) {
	atl := atlas.MustBuild(
		OpAction_AtlasEntry,
		OpActionUserinfo_AtlasEntry,
	)
	t.Run("zero opaction should roundtrip", func(t *testing.T) {
		obj := OpAction{}
		canon := `{}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := OpAction{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
	t.Run("exciting opaction should roundtrip", func(t *testing.T) {
		obj := OpAction{
			Exec:     []string{"/wizz", "bang"},
			Env:      map[string]string{"FOO": "bar", "BAZ": "quux"},
			Userinfo: &OpActionUserinfo{Username: "zoltan"},
		}
		canon := `{"exec":["/wizz","bang"],"env":{"BAZ":"quux","FOO":"bar"},"userinfo":{"username":"zoltan"}}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := OpAction{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
}

// Test Userinfo serialization.
// Some of these sanity assertions are moderately nontrivial beacuse we
// consider it important to correctly round-trip the unset/default values,
// which for integers we implement as some pointer jiggery.
func TestUserinfoSerialization(t *testing.T) {
	atl := atlas.MustBuild(OpActionUserinfo_AtlasEntry)
	t.Run("zero userinfo should roundtrip", func(t *testing.T) {
		obj := OpActionUserinfo{}
		canon := `{}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := OpActionUserinfo{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
	t.Run("userinfo with set zero uid should roundtrip", func(t *testing.T) {
		i0 := 0
		obj := OpActionUserinfo{Uid: &i0}
		canon := `{"uid":0}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := OpActionUserinfo{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, atl)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
}
