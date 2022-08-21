package rio

import (
	"github.com/polydawn/refmt/obj/atlas"
	commonatlases "github.com/polydawn/refmt/obj/atlas/common"

	api "github.com/polydawn/go-timeless-api"
)

var Atlas = atlas.MustBuild(
	Event_AtlasEntry,
	Error_AtlasEntry,
	commonatlases.Time_AsUnixInt,
	api.WareID_AtlasEntry,
)

var Event_AtlasEntry = atlas.BuildEntry((*Event)(nil)).KeyedUnion().
	Of(map[string]*atlas.AtlasEntry{
		"log":    Event_Log_AtlasEntry,
		"prog":   Event_Progress_AtlasEntry,
		"result": Event_Result_AtlasEntry,
	})

var Event_Log_AtlasEntry = atlas.BuildEntry(Event_Log{}).StructMap().Autogenerate().Complete()
var Event_Progress_AtlasEntry = atlas.BuildEntry(Event_Progress{}).StructMap().Autogenerate().Complete()
var Event_Result_AtlasEntry = atlas.BuildEntry(Event_Result{}).StructMap().Autogenerate().Complete()

var Error_AtlasEntry = atlas.BuildEntry(Error{}).StructMap().Autogenerate().Complete()

func (ll LogLevel) String() string {
	switch ll {
	case LogError:
		return "error"
	case LogWarn:
		return "warn"
	case LogInfo:
		return "info"
	case LogDebug:
		return "debug"
	default:
		return "invalid"
	}
}
