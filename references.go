package api

// n.b. nomenclature: "*Ref" is a usable tuple; "*Name" is one coordinate in a Ref.
// types are abbreviated to "*Ref"; funcs usually use the full word "*Reference*";
// funcs that are factories (or parse helpers, etc) for a type are also abbreviated.

type SlotName string
type StepName string
type SlotRef struct {
	StepName // zero for module import reference
	SlotName
}
type SubmoduleRef string // .-sep.  really is a []StepName, but we wanted something easily used as a map key.
type SubmoduleStepRef struct {
	SubmoduleRef
	StepName
}
type SubmoduleSlotRef struct {
	SubmoduleRef
	SlotRef
}

// ImportRef is a sum type, containing either a
// catalog reference ("catalog:{moduleName}:{releaseName}:{itemName}")
// or parent reference ("parent:{slotRef}"; only valid in submodules)
// or an ingest reference ("ingest:{ingestKind}[:{addntl}]"; only valid on main module).
//
// Ingest references are interesting and should be used sparingly; they're
// for where new data comes into the Timeless ecosystem -- and that also means
// ingest references are also where the Timeless Stack abilities to
// automatically recursively audit where that data came from has reached its end.
//
// Ingest references may explicitly reference wares
// (ex. "ingest:literal:tar:f00bAr"),
// or lean on other extensions to bring data into the system
// (ex. "ingest:git:.:HEAD").
// Again, use sparingly: anything beyond "ingest:literal" and your module
// pipeline has become virtually impossible for anyone to evaluate without
// whatever additional un-contained un-tracked context your ingest refers to.
//
// Ingest references should be passed on directly as an export of a module.
// Failure to do so is not *exactly* illegal, but it would make any replay
// of this module impossible without un-tracked context, and as such most of
// the tools in the Timeless Stack will issue either warnings or outright
// errors if the ingested data isn't also in the module exports.
type ImportRef interface {
	_ImportRef()
	String() string
}

type ImportRef_Catalog ItemRef
type ImportRef_Parent SlotRef
type ImportRef_Ingest struct {
	IngestKind string
	Args       string
}

func (ImportRef_Catalog) _ImportRef() {}
func (ImportRef_Parent) _ImportRef()  {}
func (ImportRef_Ingest) _ImportRef()  {}

func (x ImportRef_Catalog) String() string { return "catalog:" + (ItemRef(x)).String() }
func (x ImportRef_Parent) String() string  { return "parent:" + (SlotRef(x)).String() }
func (x ImportRef_Ingest) String() string  { return "ingest:" + "TODO" } // TODO
