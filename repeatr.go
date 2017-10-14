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
		Defines the action to perform to "evaluate" the formula -- after the
		input filesets have been assembled, these commands will be run in a
		contained sandbox on with those filesets,
		and when the commands terminate, the output filesets will be saved.

		The definition of the Action includes at minimum what commands to run,
		but also includes the option of specifying other execution parameters:
		things like environment variables, working directory, hostname...
		and (though hopefully you rarely get hung up and need to change these)
		also things like UID, GID, username, homedir, and soforth.
		All of these additional parameters have "sensible defaults" if unset.

		The Action also includes the ability to set "Policy" level -- these
		define simple privilege levels.  (The default policy is of extremely
		low privileges.)
	*/
	FormulaAction struct {
		// An array of strings to hand as args to exec -- creates a single process.
		//
		// TODO we want to add a polymorphic option here, e.g.
		// one of 'Exec', 'Script', or 'Reshuffle' may be set.
		Exec []string

		// How much power to give the process.  Default is quite low.
		Policy Policy `refmt:",omitempty"`

		// The working directory to set when invoking the executable.
		// If not set, will be defaulted to "/task".
		Cwd AbsPath `refmt:",omitempty"`

		// Environment variables.
		Env map[string]string `refmt:",omitempty"`

		// User info -- uid, gid, etc.
		Userinfo *FormulaUserinfo `refmt:",omitempty"`

		// Cradle -- enabled by default, enum value for disable.
		Cradle string `refmt:",omitempty"`

		// Hostname to set inside the container (if the executor supports this -- not all do).
		Hostname string `refmt:",omitempty"`
	}

	FormulaUserinfo struct {
		Uid      *int    `refmt:",omitempty"`
		Gid      *int    `refmt:",omitempty"`
		Username string  `refmt:",omitempty"`
		Homedir  AbsPath `refmt:",omitempty"`
	}

	OutputSpec struct {
		PackType PackType       `refmt:"packtype"`
		Filters  FilesetFilters `refmt:",omitempty"`
	}

	SetupHash string // HID of formula
)

var (
	FormulaUnion_AtlasEntry    = atlas.BuildEntry(FormulaUnion{}).StructMap().Autogenerate().Complete()
	Formula_AtlasEntry         = atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()
	FormulaContext_AtlasEntry  = atlas.BuildEntry(FormulaContext{}).StructMap().Autogenerate().Complete()
	FormulaAction_AtlasEntry   = atlas.BuildEntry(FormulaAction{}).StructMap().Autogenerate().Complete()
	FormulaUserinfo_AtlasEntry = atlas.BuildEntry(FormulaUserinfo{}).StructMap().Autogenerate().Complete()
	OutputSpec_AtlasEntry      = atlas.BuildEntry(OutputSpec{}).StructMap().Autogenerate().Complete()
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
