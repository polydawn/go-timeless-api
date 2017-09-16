package api

/*
	This file is all serializable types used in Repeatr
	to define things to run and results we get.

	WareIDs and other fileset-related basics is Rio stuff (in the 'rio.go') file.
*/

import (
	"github.com/polydawn/refmt/obj/atlas"
)

type (
	Formula struct {
		Inputs  UnpackTree
		Action  FormulaAction
		Outputs FormulaOutputs
	}

	UnpackTree map[AbsPath]UnpackSpec
	UnpackSpec struct {
		WareID  WareID         `refmt:"ware"`
		Filters FilesetFilters `refmt:"opts"`
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
	}

	FormulaOutputs map[AbsPath]string // TODO probably need more there than the ware type name ... although we could put normalizers in the "action" section

	SetupHash string // HID of formula
)

var (
	Formula_AtlasEntry       = atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()
	UnpackSpec_AtlasEntry    = atlas.BuildEntry(UnpackSpec{}).StructMap().Autogenerate().Complete()
	FormulaAction_AtlasEntry = atlas.BuildEntry(FormulaAction{}).StructMap().Autogenerate().Complete()
)

type RunRecord struct {
	UID       string             // random number, presumed globally unique.
	Time      int64              // time at start of build.
	FormulaID SetupHash          // HID of formula ran.
	Results   map[AbsPath]WareID // wares produced by the run!
	ExitCode  int                // exit code of the contained process.

	// --- below: addntl optional metadata ---

	Hostname string            // hostname.  not a trusted field, but useful for debugging.
	Metadata map[string]string // escape valve.  you can attach freetext here.
}

var RunRecord_AtlasEntry = atlas.BuildEntry(RunRecord{}).StructMap().Autogenerate().Complete()

type RunRecordHash string // HID of RunRecord.  Includes UID, etc, so quite unique.  Prefer this to UID for primary key in storage (it's collision resistant).
