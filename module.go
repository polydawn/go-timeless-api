package api

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
