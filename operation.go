package api

type AbsPath string

type Operation struct {
	Inputs  map[SlotReference]AbsPath
	Action  OpAction
	Outputs map[SlotName]AbsPath
}

type ReadyOperation struct {
	ImportsPinned map[SlotReference]WareID // any stepname here is of course opaque to repeatr, but we also don't rewrite it, so.
	Operation
}

type OpAction struct {
	Exec []string
	Env  map[string]string

	Noop bool
}
