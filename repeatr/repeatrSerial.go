package repeatr

import (
	"github.com/polydawn/refmt/obj/atlas"

	"go.polydawn.net/go-timeless-api"
)

var Atlas = atlas.MustBuild(
	Event_AtlasEntry,
	api.OperationRecord_AtlasEntry,
	api.WareID_AtlasEntry,
)

var Event_AtlasEntry = atlas.BuildEntry((*Event)(nil)).KeyedUnion().
	Of(map[string]*atlas.AtlasEntry{
		"log":    Event_Log_AtlasEntry,
		"txt":    Event_Output_AtlasEntry,
		"result": Event_Result_AtlasEntry,
	})

var Event_Log_AtlasEntry = atlas.BuildEntry(Event_Log{}).StructMap().Autogenerate().Complete()
var Event_Output_AtlasEntry = atlas.BuildEntry(Event_Output{}).StructMap().Autogenerate().Complete()
var Event_Result_AtlasEntry = atlas.BuildEntry(Event_Result{}).StructMap().Autogenerate().Complete()
