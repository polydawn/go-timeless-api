package api

import (
	"github.com/polydawn/refmt/obj/atlas"
)

var rioAtlasEntries = []*atlas.AtlasEntry{
	WareID_AtlasEntry,
	FilesetFilters_AtlasEntry,
}

var repeatrAtlasEntries = []*atlas.AtlasEntry{
	FormulaUnion_AtlasEntry,
	Formula_AtlasEntry,
	FormulaContext_AtlasEntry,
	FormulaAction_AtlasEntry,
	FormulaUserinfo_AtlasEntry,
	OutputSpec_AtlasEntry,
	RunRecord_AtlasEntry,
}

var formulaCasAtlasEntries = []*atlas.AtlasEntry{
	Formula_AtlasEntry,
	FormulaAction_AtlasEntry,
	FormulaUserinfo_AtlasEntry,
	OutputSpec_AtlasEntry,
}

var hitchAtlasEntries = []*atlas.AtlasEntry{
	Catalog_AtlasEntry,
	ReleaseItemID_AtlasEntry,
	ReleaseEntry_AtlasEntry,
	Replay_AtlasEntry, // probably delete now
	Step_AtlasEntry,   // probably delete now
	Basting_AtlasEntry,
	BastingStep_AtlasEntry,
}

var RepeatrAtlas = atlas.MustBuild(
	aecat(
		rioAtlasEntries,
		repeatrAtlasEntries,
	)...,
)

var FormulaCasAtlas = atlas.MustBuild(
	aecat(
		rioAtlasEntries,
		formulaCasAtlasEntries,
	)...,
)

var HitchAtlas = atlas.MustBuild(
	aecat(
		rioAtlasEntries,
		repeatrAtlasEntries,
		hitchAtlasEntries,
	)...,
)

func aecat(aess ...[]*atlas.AtlasEntry) (r []*atlas.AtlasEntry) {
	for _, aes := range aess {
		r = append(r, aes...)
	}
	return
}
