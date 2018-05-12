package api

type SlotName string
type StepName string
type StepReference []string // .-sep
type SlotReference struct {
	StepName // zero for module import reference
	SlotName
}

type ImportPattern []string   // :-sep.  Major types catalog|ingest|ware|parent
type ImportReference []string // :-sep.  Major types catalog|ware
type WareID [2]string         // :-sep

type ItemName string // 3rd part of catalog index tuple

type Module struct {
	Imports         map[SlotName]ImportPattern
	ImportsResolved map[SlotName]ImportReference // always nil on submodules // FIXME key may be a DEEP slot reference
	ImportsPinned   map[SlotName]WareID          // always nil on submodules // FIXME key may be a DEEP slot reference // also you can't have slices in map keys lol
	Operation       *StepUnion
	Operations      map[StepName]StepUnion
	Exports         map[ItemName]SlotReference
}

type StepUnion interface {
	_Step()
}

func (Module) _Step()    {}
func (Operation) _Step() {}
