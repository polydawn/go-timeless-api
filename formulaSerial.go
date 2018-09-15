package api

import "github.com/polydawn/refmt/obj/atlas"

var (
	Atlas_Formula = atlas.MustBuild(
		Formula_AtlasEntry,
		FormulaAction_AtlasEntry,
		FormulaUserinfo_AtlasEntry,
		FormulaOutputSpec_AtlasEntry,
		FilesetPackFilter_AtlasEntry,
		WareID_AtlasEntry,
	)

	Atlas_FormulaRunRecord = atlas.MustBuild(
		FormulaRunRecord_AtlasEntry,
		WareID_AtlasEntry,
	)
)

var (
	Formula_AtlasEntry           = atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()
	FormulaAction_AtlasEntry     = atlas.BuildEntry(FormulaAction{}).StructMap().Autogenerate().Complete()
	FormulaUserinfo_AtlasEntry   = atlas.BuildEntry(FormulaUserinfo{}).StructMap().Autogenerate().Complete()
	FormulaOutputSpec_AtlasEntry = atlas.BuildEntry(FormulaOutputSpec{}).StructMap().Autogenerate().Complete()
	FormulaRunRecord_AtlasEntry  = atlas.BuildEntry(FormulaRunRecord{}).StructMap().Autogenerate().Complete()
)
