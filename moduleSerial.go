package api

import (
	"github.com/polydawn/refmt/obj/atlas"
)

var Atlas_Module = atlas.MustBuild(
	Module_AtlasEntry,
	StepUnion_AtlasEntry,
	Operation_AtlasEntry,
	SlotRef_AtlasEntry,
	ImportRef_AtlasEntry,
	FormulaAction_AtlasEntry,
	FormulaUserinfo_AtlasEntry,
)

var Module_AtlasEntry = atlas.BuildEntry(Module{}).StructMap().Autogenerate().Complete()

var StepUnion_AtlasEntry = atlas.BuildEntry((*StepUnion)(nil)).KeyedUnion().
	Of(map[string]*atlas.AtlasEntry{
		"module":    Module_AtlasEntry,
		"operation": Operation_AtlasEntry,
	})
