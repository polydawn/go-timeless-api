package api

type SlotName string
type StepName string
type SlotReference struct {
	StepName // zero for module import reference
	SlotName
}
type SubmoduleReference string // .-sep.  really is a []StepName, but we wanted something easily used as a map key.
type SubmoduleSlotReference struct {
	SubmoduleReference
	SlotName
}

type ImportPattern []string   // :-sep.  Major types catalog|ingest|ware|parent
type ImportReference []string // :-sep.  Major types catalog|ware

type ItemName string // 3rd part of catalog index tuple

type Module struct {
	Imports         map[SlotName]ImportPattern
	ImportsResolved map[SubmoduleSlotReference]ImportReference // only allowed on top module (since it contains info for all submodules as well).
	ImportsPinned   map[SubmoduleSlotReference]WareID          // only allowed on top module (since it contains info for all submodules as well).
	Operation       *StepUnion
	Operations      map[StepName]StepUnion
	Exports         map[ItemName]SlotReference
}

type StepUnion interface {
	_Step()
}

func (Module) _Step()    {}
func (Operation) _Step() {}
