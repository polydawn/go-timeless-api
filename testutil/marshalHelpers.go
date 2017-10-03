package testutil

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	"github.com/polydawn/refmt/obj/atlas"
)

func ShouldMarshalJson(t *testing.T, thing interface{}, atl atlas.Atlas, fixture string) {
	t.Helper()
	msg, err := refmt.MarshalAtlased(
		json.EncodeOptions{},
		thing,
		atl,
	)
	AssertNoError(t, err)
	WantStringEqual(t, string(msg), fixture)
}

func ShouldMarshalPrettyJson(t *testing.T, thing interface{}, atl atlas.Atlas, fixture string) {
	t.Helper()
	msg, err := refmt.MarshalAtlased(
		json.EncodeOptions{},
		thing,
		atl,
	)
	AssertNoError(t, err)
	WantStringEqual(t, PrettifyJson(msg), Dedent(fixture))
}

func ShouldUnmarshalJson(t *testing.T, thing string, atl atlas.Atlas, matching interface{}) {
	t.Helper()
	slot := reflect.New(reflect.ValueOf(matching).Type()).Interface()
	err := refmt.UnmarshalAtlased(
		json.DecodeOptions{},
		[]byte(thing),
		slot,
		atl,
	)
	AssertNoError(t, err)
	diff := cmp.Diff(reflect.ValueOf(slot).Elem().Interface(), matching)
	if diff != "" {
		t.Errorf("Match failed: struct diff:\n%s", diff)
	}
}
