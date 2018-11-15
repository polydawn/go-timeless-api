package api

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
