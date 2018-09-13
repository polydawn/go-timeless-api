package rio

import (
	"github.com/polydawn/refmt/obj/atlas"

	"go.polydawn.net/go-timeless-api"
)

var Atlas = atlas.MustBuild(
	Event_AtlasEntry,
	api.WareID_AtlasEntry,
)

var Event_AtlasEntry = atlas.BuildEntry((*Event)(nil)).KeyedUnion().
	Of(map[string]*atlas.AtlasEntry{
		"log":    Event_Log_AtlasEntry,
		"porg":   Event_Progress_AtlasEntry,
		"result": Event_Result_AtlasEntry,
	})

var Event_Log_AtlasEntry = atlas.BuildEntry(Event_Log{}).StructMap().Autogenerate().Complete()
var Event_Progress_AtlasEntry = atlas.BuildEntry(Event_Progress{}).StructMap().Autogenerate().Complete()
var Event_Result_AtlasEntry = atlas.BuildEntry(Event_Result{}).StructMap().Autogenerate().Complete()

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
