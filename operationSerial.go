package api

import "github.com/polydawn/refmt/obj/atlas"

var (
	Operation_AtlasEntry       = atlas.BuildEntry(Operation{}).StructMap().Autogenerate().Complete()
	OperationRecord_AtlasEntry = atlas.BuildEntry(OperationRecord{}).StructMap().Autogenerate().Complete()
)
