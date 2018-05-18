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

type ImportPattern []string   // :-sep.  Major types catalog|ingest|ware|parent
type ImportReference []string // :-sep.  Major types catalog|ware

type ItemName string // 3rd part of catalog index tuple

type Module struct {
	Imports         map[SlotName]ImportPattern
	ImportsResolved map[SubmoduleSlotReference]ImportReference // only allowed on top module (since it contains info for all submodules as well).
	ImportsPinned   map[SubmoduleSlotReference]WareID          // only allowed on top module (since it contains info for all submodules as well).
	Operations      map[StepName]StepUnion
	Exports         map[ItemName]SlotReference
}

type StepUnion interface {
	_Step()
}

func (Module) _Step()    {}
func (Operation) _Step() {}

func submodularizeReference(parent StepName, ref SubmoduleReference) SubmoduleReference {
	if ref == "" {
		return SubmoduleReference(parent)
	}
	return SubmoduleReference(string(parent) + "." + string(ref))
}
func submodularizeStepReference(parent StepName, ref SubmoduleStepReference) SubmoduleStepReference {
	return SubmoduleStepReference{
		submodularizeReference(parent, ref.SubmoduleReference),
		ref.StepName,
	}
}

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
