package api

import (
	"github.com/polydawn/refmt/obj/atlas"
)

type Basting struct {
	Steps    map[string]BastingStep
	Contexts map[string]FormulaContext `refmt:",omitempty"`
}

type BastingStep struct {
	// Named imports for all the formula inputs.
	// These may be either "{catalog}:{version}:{item}" tuples, or
	// the basting-local "wire:{step}:{output}" tuple.
	Imports map[AbsPath]ReleaseItemID

	// The formula to run for this step.
	// The 'action' and 'outputs' sections should certainly be complete;
	// the 'input' section *may* be missing its hashes (definitely blank for
	// "wire" imports, which cannot be filled in until we're executing the
	// whole group of basted steps; possibly for named catalog imports, which
	// can be resolved at any time by referring to a hitch database).
	Formula Formula

	// no FormulaContext -- that's joined only right before 'repeatr run'.
	// You can find them in the Basting with matching step names -- but only in some states:
	//  - Basting that's about to be executed should have the context URLs provided.
	//  - Basting that's stored by e.g. hitch releases should not.
}

var (
	Basting_AtlasEntry     = atlas.BuildEntry(Basting{}).StructMap().Autogenerate().Complete()
	BastingStep_AtlasEntry = atlas.BuildEntry(BastingStep{}).StructMap().Autogenerate().Complete()
)
