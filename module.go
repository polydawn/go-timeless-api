package api

type SlotName string
type StepName string
type SlotReference struct {
	StepName // zero for module import reference
	SlotName
}
type SubmoduleReference string // .-sep.  really is a []StepName, but we wanted something easily used as a map key.
type SubmoduleStepReference struct {
	SubmoduleReference
	StepName
}
type SubmoduleSlotReference struct {
	SubmoduleReference
	SlotReference
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
}

type ImportRef_Catalog ItemRef
type ImportRef_Parent SlotReference
type ImportRef_Ingest []string

func (ImportRef_Catalog) _ImportRef() {}
func (ImportRef_Parent) _ImportRef()  {}
func (ImportRef_Ingest) _ImportRef()  {}

type MainModule struct {
	ImportsPinned map[SubmoduleSlotReference]WareID // only allowed on top module (since it contains info for all submodules as well).
	Module
}
type Module struct {
	Imports    map[SlotName]ImportRef
	Operations map[StepName]StepUnion
	Exports    map[ItemName]SlotReference
}

type StepUnion interface {
	_Step()
}

func (Module) _Step()    {}
func (Operation) _Step() {}

func validateConnectivity(m Module) ([]StepName, []error) {
	// Suppose all imports are unused; we'll strike things off as they're used.
	unusedImports := make(map[SlotName]struct{}, len(m.Imports))
	for imp := range m.Imports {
		unusedImports[imp] = struct{}{}
	}

	return nil, nil
}

func validateParentImports(parent Module, submodule Module) []error {
	return nil
}