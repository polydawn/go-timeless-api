package api

import (
	"github.com/polydawn/refmt/obj/atlas"
)

var Atlas_Module = atlas.MustBuild(
	Module_AtlasEntry,
	Operation_AtlasEntry,
	SlotReference_AtlasEntry,
	OpAction_AtlasEntry,
	OpActionUserinfo_AtlasEntry,
)

var Module_AtlasEntry = atlas.BuildEntry(Module{}).StructMap().Autogenerate().Complete()
