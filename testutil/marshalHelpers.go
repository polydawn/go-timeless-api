package testutil

import (
	"testing"

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
