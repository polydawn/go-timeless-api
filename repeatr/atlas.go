package repeatr

import (
	"github.com/polydawn/refmt/obj/atlas"
	"github.com/polydawn/refmt/obj/atlas/common"

	"go.polydawn.net/go-timeless-api"
)

var Atlas = atlas.MustBuild(
	atlas.BuildEntry(Event{}).StructMap().Autogenerate().Complete(),
	atlas.BuildEntry(Event_Log{}).StructMap().Autogenerate().Complete(),
	atlas.BuildEntry(Event_Output{}).StructMap().Autogenerate().Complete(),
	atlas.BuildEntry(Event_Result{}).StructMap().Autogenerate().Complete(),
	atlas.BuildEntry(Error{}).StructMap().Autogenerate().Complete(),
	commonatlases.Time_AsUnixInt,
	api.RunRecord_AtlasEntry,
	api.WareID_AtlasEntry,
)
