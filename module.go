package api

type MainModule struct {
	ImportsPinned map[SubmoduleSlotRef]WareID // only allowed on top module (since it contains info for all submodules as well).
	Module
}
type Module struct {
	Imports map[SlotName]ImportRef
	Steps   map[StepName]StepUnion
	Exports map[ItemName]SlotRef `refmt:",omitempty"`
}

type StepUnion interface {
	_Step()
}

func (Module) _Step()    {}
func (Operation) _Step() {}
