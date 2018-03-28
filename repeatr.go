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

		  - The Formula itself is always used.
		  - The Formula must have all input hashes resolved if it's about
		    to be evaluated by `repeatr run`.
		  - The Imports field should be nil if we're about to `repeatr run`,
		    because Repeatr itself does not know or care about higher level graphs,
		    but we won't explicitly reject it either.  (If there are import
		    paths without a corresponding input hash, though, that's something
		    gone wrong somewhere and we'll halt and warn.)
		  - The Context fields should be set if we're about to `repeatr run`.

		Note that the Basting structure contains all of these same pieces of
		info, but composed differently: this is because at every level of
		composition, we want the "context" (URLs and such) to remain separate
		from the hashable, sharable, deterministic, timeless parts.

		Some code uses this same union in marshalling partial info during
		generation of Basting.  In that case:

		  - If being used in Basting, the Imports fields may be assigned,
		    and some input hashes may be blank if they're "wire" imports.
		  - The Context fields may also be used temporarily (but as already
		    mentioned, this will be destructured when composed into the Basting).
	*/
	FormulaUnion struct {
		Imports map[AbsPath]ReleaseItemID `refmt:",omitempty"`
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

/*
	Much as per the Formula.Apply function, but also combines the Imports
	and the Context if present.
*/
func (f1 FormulaUnion) Apply(f2 FormulaUnion) (f3 FormulaUnion) {
	refmt.MustCloneAtlased(f1.Imports, &f3.Imports, RepeatrAtlas)
	if f2.Imports != nil {
		if f3.Imports == nil {
			f3.Imports = make(map[AbsPath]ReleaseItemID)
		}
		for p, id := range f2.Imports {
			f3.Imports[p] = id
		}
	}

	f3.Formula = f1.Formula.Apply(f2.Formula)

	refmt.MustCloneAtlased(f1.Context, &f3.Context, RepeatrAtlas)
	if f2.Context != nil {
		if f3.Context == nil {
			f3.Context = &FormulaContext{}
		}
		if f2.Context.FetchUrls != nil {
			if f3.Context.FetchUrls == nil {
				f3.Context.FetchUrls = make(map[AbsPath][]WarehouseAddr)
			}
			for p, addrs := range f2.Context.FetchUrls {
				f3.Context.FetchUrls[p] = append(f3.Context.FetchUrls[p], addrs...)
			}
		}
		if f2.Context.SaveUrls != nil {
			if f3.Context.SaveUrls == nil {
				f3.Context.SaveUrls = make(map[AbsPath]WarehouseAddr)
			}
			f3.Context.SaveUrls = f3.Context.SaveUrls
		}
	}
	return
}

func (f *Formula) Clone() (f2 Formula) {
	refmt.MustCloneAtlased(f, &f2, RepeatrAtlas)
	return
}

/*
	Apply another formula as a patch to this one, returning a new formula
	with the combined content.

	Fields that are primitives will take the value from f2 if set.
	Fields that are maps -- such as inputs, outputs, and action.env -- are
	joined, and any duplicate keys in f2 will override the values from f1.
	Fields that are arrays will either be merged by appending values from f2,
	or, in some cases, may be treated as if they are primitives (action.exec
	is treated as a primitive, for example, and if set at all in f2, will
	override any values from f1 completely).

*/
func (f1 Formula) Apply(f2 Formula) (f3 Formula) {
	f3 = f1.Clone()
	if f3.Inputs == nil {
		f3.Inputs = make(map[AbsPath]WareID)
	}
	for k, v := range f2.Inputs {
		f3.Inputs[k] = v
	}
	{
		if f2.Action.Exec != nil {
			f3.Action.Exec = f2.Action.Exec
		}
		if f2.Action.Policy != "" {
			f3.Action.Policy = f2.Action.Policy
		}
		if f2.Action.Cwd != "" {
			f3.Action.Cwd = f2.Action.Cwd
		}
		if f3.Action.Env == nil {
			f3.Action.Env = make(map[string]string)
		}
		for k, v := range f2.Action.Env {
			f3.Action.Env[k] = v
		}
		if f2.Action.Userinfo != nil {
			if f3.Action.Userinfo == nil {
				f3.Action.Userinfo = &FormulaUserinfo{}
			}
			if f2.Action.Userinfo.Uid != nil {
				i := *f2.Action.Userinfo.Uid
				f3.Action.Userinfo.Uid = &i
			}
			if f2.Action.Userinfo.Gid != nil {
				i := *f2.Action.Userinfo.Gid
				f3.Action.Userinfo.Gid = &i
			}
			if f2.Action.Userinfo.Username != "" {
				f3.Action.Userinfo.Username = f2.Action.Userinfo.Username
			}
			if f2.Action.Userinfo.Homedir != "" {
				f3.Action.Userinfo.Homedir = f2.Action.Userinfo.Homedir
			}
		}
		if f2.Action.Cradle != "" {
			f3.Action.Cradle = f2.Action.Cradle
		}
		if f2.Action.Hostname != "" {
			f3.Action.Hostname = f2.Action.Hostname
		}
	}
	if f3.Outputs == nil {
		f3.Outputs = make(map[AbsPath]OutputSpec)
	}
	for k, v := range f2.Outputs {
		f3.Outputs[k] = v
	}
	return
}

type RunRecord struct {
	Guid      string             // random number, presumed globally unique.
	Time      int64              // time at start of build.
	FormulaID SetupHash          // HID of formula ran.
	ExitCode  int                // exit code of the contained process.
	Results   map[AbsPath]WareID // wares produced by the run!

	// --- below: addntl optional metadata ---

	Hostname string            `refmt:",omitempty"` // hostname.  not a trusted field, but useful for debugging.
	Metadata map[string]string `refmt:",omitempty"` // escape valve.  you can attach freetext here.
}

var RunRecord_AtlasEntry = atlas.BuildEntry(RunRecord{}).StructMap().Autogenerate().Complete()

type RunRecordHash string // HID of RunRecord.  Includes guid, etc, so quite unique.  Prefer this to guid for primary key in storage (it's collision resistant).
