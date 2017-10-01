package api

/*
	This file is all serializable types used in Repeatr
	to define things to run and results we get.

	WareIDs and other fileset-related basics is Rio stuff (in the 'rio.go') file.
*/

import (
	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/obj/atlas"
)

type (
	/*
		FormulaUnion is a serialization helper struct which composes several
		structures that are typically saved in the same file for convenience.
		FormulaUnion is not used in any function APIs -- component fields are.
	*/
	FormulaUnion struct {
		Formula Formula
		Context *FormulaContext `refmt:",omitempty"`
	}

	/*
		The essense of Repeatr: a Formula is full instructions for populating
		a sandbox and running a process in it, then collecting and saving
		output filesets.

		The formula, serialized and hashed, produces a SetupHash.
		This SetupHash string is effectively a memoization key for the
		entire sandboxed computation.
	*/
	Formula struct {
		Inputs  map[AbsPath]WareID
		Action  FormulaAction
		Outputs map[AbsPath]OutputSpec
	}

	/*
		FormulaContext contains the remaining references that are needed to
		actually evaluate a Formula.  Locations for fetching inputs and
		storing outputs, etc.

		FormulaContext is separated out from the Formula because this
		information is *not* part of the SetupHash -- it's not relevant to
		the reproducibility of the setup or memoization of the computation.
		They're often serialized in the same file though (see FormulaUnion).
	*/
	FormulaContext struct {
		FetchUrls map[AbsPath][]WarehouseAddr
		SaveUrls  map[AbsPath]WarehouseAddr
	}

	/*
		Defines the action to perform to evaluate the formula -- some commands
		or filesystem operations which will be run after the inputs have been
		assembled; the action is done, the outputs will be saved.
	*/
	FormulaAction struct {
		// An array of strings to hand as args to exec -- creates a single process.
		//
		// TODO we want to add a polymorphic option here, e.g.
		// one of 'Exec', 'Script', or 'Reshuffle' may be set.
		Exec []string

		// The working directory to set when invoking the executable.
		// If not set, will be defaulted to "/task".
		Cwd AbsPath `refmt:",omitempty"`

		// Environment variables.
		Env map[string]string `refmt:",omitempty"`

		// Hostname to set inside the container (if the executor supports this -- not all do).
		Hostname string `refmt:",omitempty"`
	}

	OutputSpec struct {
		PackType PackType       `refmt:"packtype"`
		Filters  FilesetFilters `refmt:",omitempty"`
	}

	SetupHash string // HID of formula
)

var (
	FormulaUnion_AtlasEntry   = atlas.BuildEntry(FormulaUnion{}).StructMap().Autogenerate().Complete()
	Formula_AtlasEntry        = atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()
	FormulaContext_AtlasEntry = atlas.BuildEntry(FormulaContext{}).StructMap().Autogenerate().Complete()
	FormulaAction_AtlasEntry  = atlas.BuildEntry(FormulaAction{}).StructMap().Autogenerate().Complete()
	OutputSpec_AtlasEntry     = atlas.BuildEntry(OutputSpec{}).StructMap().Autogenerate().Complete()
)

func (f *Formula) Clone() (f2 Formula) {
	refmt.MustCloneAtlased(f, &f2, RepeatrAtlas)
	return
}

type RunRecord struct {
	Guid      string             // random number, presumed globally unique.
	Time      int64              // time at start of build.
	FormulaID SetupHash          // HID of formula ran.
	Results   map[AbsPath]WareID // wares produced by the run!
	ExitCode  int                // exit code of the contained process.

	// --- below: addntl optional metadata ---

	Hostname string            // hostname.  not a trusted field, but useful for debugging.
	Metadata map[string]string // escape valve.  you can attach freetext here.
}

var RunRecord_AtlasEntry = atlas.BuildEntry(RunRecord{}).StructMap().Autogenerate().Complete()

type RunRecordHash string // HID of RunRecord.  Includes guid, etc, so quite unique.  Prefer this to guid for primary key in storage (it's collision resistant).
