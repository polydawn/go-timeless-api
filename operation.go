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

/*
	OperationRecord is mostly an alias of FormulaRunRecord, but with Results
	indexed by SlotName from the Operation rather than path in the Formula.

	We usually serialize FormulaRunRecord, because it's more convergent when
	content-addressed; OperationRecord contains immaterial details (e.g. the
	SlotName).  OperationRecord is sometimes more convenient to use internally.
*/
type OperationRecord struct {
	FormulaRunRecord
	Results map[SlotName]WareID
}
