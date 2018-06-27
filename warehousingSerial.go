package api

import "github.com/polydawn/refmt/obj/atlas"

var (
	WareSourcing_AtlasEntry = atlas.BuildEntry(WareSourcing{}).StructMap().Autogenerate().Complete()
)

var Atlas_WareSourcing = atlas.MustBuild(
	WareSourcing_AtlasEntry,
	WareID_AtlasEntry,
)
