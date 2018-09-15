package api

import "github.com/polydawn/refmt/obj/atlas"

var (
	Formula_AtlasEntry           = atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()
	FormulaOutputSpec_AtlasEntry = atlas.BuildEntry(FormulaOutputSpec{}).StructMap().Autogenerate().Complete()
	FormulaAction_AtlasEntry     = atlas.BuildEntry(FormulaAction{}).StructMap().Autogenerate().Complete()
	FormulaUserinfo_AtlasEntry   = atlas.BuildEntry(FormulaUserinfo{}).StructMap().Autogenerate().Complete()
)
