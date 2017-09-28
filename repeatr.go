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
	Formula struct {
		Inputs    map[AbsPath]WareID
		Action    FormulaAction
		Outputs   map[AbsPath]OutputSpec
		FetchUrls map[AbsPath][]WarehouseAddr
		SaveUrls  map[AbsPath][]WarehouseAddr
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

	OutputSpec struct {
		PackFmt string         `refmt:"packfmt"`
		Filters FilesetFilters `refmt:",omitempty"`
	}

	SetupHash string // HID of formula
)

var (
	Formula_AtlasEntry    = atlas.BuildEntry(Formula{}).StructMap().Autogenerate().Complete()
	FormulaCas_AtlasEntry = atlas.BuildEntry(Formula{}).StructMap().
				AddField("Inputs", atlas.StructMapEntry{SerialName: "inputs"}).
				AddField("Action", atlas.StructMapEntry{SerialName: "action"}).
				AddField("Outputs", atlas.StructMapEntry{SerialName: "outputs"}).
				Complete() // Note the explicit lack of fetchUrls and saveUrls.
	FormulaAction_AtlasEntry = atlas.BuildEntry(FormulaAction{}).StructMap().Autogenerate().Complete()
	OutputSpec_AtlasEntry    = atlas.BuildEntry(OutputSpec{}).StructMap().Autogenerate().Complete()
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
