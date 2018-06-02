package api

import "github.com/polydawn/refmt/obj/atlas"

var (
	OpAction_AtlasEntry         = atlas.BuildEntry(OpAction{}).StructMap().Autogenerate().Complete()
	OpActionUserinfo_AtlasEntry = atlas.BuildEntry(OpActionUserinfo{}).StructMap().Autogenerate().Complete()
)
