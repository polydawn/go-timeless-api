package api

type AbsPath string

/*
	Operation is one of the concrete types of StepUnion which composes a Module;
	it describes a containerizable computation, all of its input filesystem
	paths bound to slot references, and all of the paths that should be collected
	as outputs and assigned to another slot for further use.

	When all of the input slot references in an Operation are known, it can be
	bound, becoming a Formula -- which is structurally similar, but now with
	all specific, concrete WareID hashes instead of SlotRef.
*/
type Operation struct {
	Inputs  map[SlotRef]AbsPath
	Action  FormulaAction
	Outputs map[SlotName]AbsPath `refmt:",omitempty"`
}

type OperationRecord struct {
	Guid     string              // random number, presumed globally unique.
	Time     int64               // time at start of build.
	ExitCode int                 // exit code of the contained process.
	Results  map[SlotName]WareID // wares produced by the run!
}
