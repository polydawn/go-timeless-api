package repeatr

import (
	"context"
	"fmt"

	"go.polydawn.net/go-timeless-api"
)

// RunOperation evaluates an api.Operation, using a repeatr.RunFunc.
//
// RunOperation takes care of all the details of converting human-readable
// slot names in the operation into bound WareIDs, and back again for output
// results, etc.
//
// Formulas index everything by paths, because paths are material to the
// outcome of the computation, and slot names as used in Operations *are not*.
// This makes Formulas, when hashed, a useful "primary key" for other lookups;
// but it means we need to pivot some things to go from Operation to Formula.
//
// The warehousing parameters of this function are similarly more opinionated
// than the level of detail that the RunFunc interface allowes; since
// RunOperation is intended to be used when doing lots of operations, it makes
// more sense to use WareSourcing and WareStaging configuration, and steer
// toward output warehouses that are the most reusable: only content
// addressable warehouses, indexed by packType, are allowed arguments.
//
// Note that the Monitor system is not translated, and also uses paths rather
// than slotrefs and slotnames, and still reports results in FormulaRunRecord.
func RunOperation(
	ctx context.Context,
	runTool RunFunc, // the repeatr API to drive
	op api.Operation, // What protype of a formula to bind and run.
	scope map[api.SlotRef]api.WareID, // What slots are in scope to reference as inputs.
	wareSourcing api.WareSourcing, // Suggestions on where to get wares.
	wareStaging api.WareStaging, // Instructions on where to store output wares.
	input InputControl, // Optionally: input control.  The zero struct is no input (which is fine).
	monitor Monitor, // Optionally: callbacks for progress monitoring.  Also where stdout/stderr is gathered.
) (*api.OperationRecord, error) {
	// Convert the Operation into a fully bound Formula.
	//  This has several steps:
	//    - Look up the input wareIDs from the 'scope' map;
	//    - Organize those by input paths;
	//    - and convert the output specs to also be indexed by path.
	//  (We'll flip warehousing separately after this.)
	//  We'll hang on to an 'outputReverseMap' at the end, which we need later
	//   in order to flip the results back to their original named references.
	frm := api.Formula{
		Inputs:  make(map[api.AbsPath]api.WareID, len(op.Inputs)),
		Action:  op.Action,
		Outputs: make(map[api.AbsPath]api.FormulaOutputSpec, len(op.Outputs)),
	}
	for slotRef, pth := range op.Inputs {
		pin, ok := scope[slotRef]
		if !ok {
			return nil, fmt.Errorf("cannot provide an input ware for slotref %q: no such ref in scope", slotRef)
		}
		frm.Inputs[pth] = pin
	}
	outputReverseMap := make(map[api.AbsPath]api.SlotName, len(op.Outputs))
	for slotName, pth := range op.Outputs {
		frm.Outputs[pth] = api.FormulaOutputSpec{
			PackType: "tar", // FUTURE: the api.Operation layer doesn't really support config for this yet, but should
			// FUTURE: filters are also conspicuously missing at the api.Operation layer, and should be included
		}
		outputReverseMap[pth] = slotName
	}

	// Traverse the Formula and populate FormulaContext from WareSourcing
	//  and plugging output storage according to WareStaging.
	wareSourcing = wareSourcing.PivotToInputs(frm)
	frmCtx := FormulaContext{
		FetchUrls: make(map[api.AbsPath][]api.WarehouseLocation),
		SaveUrls:  make(map[api.AbsPath]api.WarehouseLocation),
	}
	for pth, wareID := range frm.Inputs {
		frmCtx.FetchUrls[pth] = wareSourcing.ByWare[wareID]
	}
	for pth, outSpec := range frm.Outputs {
		frmCtx.SaveUrls[pth] = wareStaging.ByPackType[outSpec.PackType]
	}

	// All flipped.  Run!
	runRecord, err := runTool(
		ctx,
		frm,
		frmCtx,
		input,
		monitor,
	)
	if err != nil {
		return nil, err
	}

	// Convert outputs back to indexed by the slotnames we started with.
	opRecord := api.OperationRecord{
		FormulaRunRecord: *runRecord,
		Results:          make(map[api.SlotName]api.WareID, len(runRecord.Results)),
	}
	for pth, wareID := range runRecord.Results {
		opRecord.Results[outputReverseMap[pth]] = wareID
	}

	// Ya ta!
	return &opRecord, nil
}
