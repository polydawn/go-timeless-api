package api

import "github.com/polydawn/refmt/obj/atlas"

var (
	Operation_AtlasEntry        = atlas.BuildEntry(Operation{}).StructMap().Autogenerate().Complete()
	OpAction_AtlasEntry         = atlas.BuildEntry(OpAction{}).StructMap().Autogenerate().Complete()
	OpActionUserinfo_AtlasEntry = atlas.BuildEntry(OpActionUserinfo{}).StructMap().Autogenerate().Complete()
)
