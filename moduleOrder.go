package api

import (
	"fmt"
	"sort"
)

/*
	OrderStepsDeep is like OrderSteps, but returns *all* steps, recursively
	including all submodules and their steps.
*/
func OrderStepsDeep(m Module) (r []SubmoduleStepReference, _ error) {
	levelOrder, err := OrderSteps(m)
	if err != nil {
		return nil, err
	}
	for _, stepName := range levelOrder {
		switch x := m.Operations[stepName].(type) {
		case Operation:
			r = append(r, SubmoduleStepReference{"", stepName})
		case Module:
			subOrder, err := OrderStepsDeep(x)
			if err != nil {
				return nil, err
			}
			for _, subStep := range subOrder {
				r = append(r, subStep.Contextualize(stepName))
			}
		}
	}
	return r, nil
}

/*
	OrderSteps returns a slice of StepName from a Module's Operations
	in order of topological sort based on SlotName each op references as inputs.

	We break ties based on lexigraphical sort on the step names.
	We choose this simple tie-breaker rather than attempting any fancier
	logic based on e.g. downstream dependencies, etc, because ease of
	understanding and the simplicity of predicting the result of the sort
	is more important than cleverness; so is the regional stability of the
	sort in the face of changes in other parts of the graph.
*/
func OrderSteps(m Module) ([]StepName, error) {
	// Alloc result accumulator.
	result := make([]StepName, 0, len(m.Operations))
	// Initialize todo set; it shrinks as we go.
	todo := make(map[StepName]struct{}, len(m.Operations))
	for step := range m.Operations {
		todo[step] = struct{}{}
	}
	// Sort operations by their name (this is our tiebreaker, in advance).
	stepsOrdered := make(stepNames, 0, len(m.Operations))
	for step := range m.Operations {
		stepsOrdered = append(stepsOrdered, step)
	}
	sort.Sort(stepsOrdered)
	// For each step: visit.  (This will recurse, and no-op itself internally as approrpriate for visited nodes.)
	for _, step := range stepsOrdered {
		if err := orderSteps_visit(step, todo, map[StepName]struct{}{}, &result, m); err != nil {
			return nil, err
		}
	}
	return result, nil
}

type stepNames []StepName

func (s stepNames) Len() int           { return len(s) }
func (s stepNames) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s stepNames) Less(i, j int) bool { return s[i] < s[j] }

// func OrderStepsYielding(m Module, want []SlotReference) ([]StepName, error)
//   Easier to implement this later by reverse walking the toposort of []StepName
//    that we've already computed and simply dropping irrelevant ones.
//    This approach has a lower CPU time if we've already done the toposort;
//    and we *always* want to have already done the toposort, because in final
//    UX we would always want our tool to return errors on invalid topology
//    even if when running in the "just this swatch of module" mode you
//    coincidentally wouldn't happen to hit the invalid region.

func orderSteps_visit(
	node StepName,
	todo map[StepName]struct{},
	loopDetector map[StepName]struct{},
	result *[]StepName,
	m Module,
) error {
	// Quick exit if possible.
	if _, ok := todo[node]; !ok {
		return nil
	}
	if _, ok := loopDetector[node]; ok {
		return fmt.Errorf("not a dag: loop detected at %q", node)
	}
	// Mark self for loop detection.
	loopDetector[node] = struct{}{}
	// Extract any imports which are dependency wiring.
	wires := inputSlotReferences(m.Operations[node])
	// Check that those actually point somewhere.
	for _, wire := range wires {
		// TODO: all of these name existence checks should be done linearly up front (... also).
		// FIXME: these error strings are not great for submodules; they forgot the "parent:" part.
		switch wire.StepName == "" {
		case true:
			if _, ok := m.Imports[wire.SlotName]; !ok {
				return fmt.Errorf("step %q has an invalid reference to %q: %q is not the name of an import in this module", node, wire, wire.SlotName)
			}
		case false:
			if op, ok := m.Operations[wire.StepName]; !ok {
				return fmt.Errorf("step %q has an invalid reference to %q: %q is not the name of a step in this module", node, wire, wire.SlotName)
			} else if _, ok := outputSlotReferences(op)[wire.SlotName]; !ok {
				return fmt.Errorf("step %q has an invalid reference to %q: step %q has no output named %s", node, wire, wire.StepName, wire.SlotName)
			}
		}
	}
	// Sort the dependency nodes by name, then recurse.
	//  This sort is necessary for deterministic order of unrelated nodes.
	sort.Sort(slotRefs(wires))
	for _, wire := range wires {
		switch wire.StepName == "" {
		case true:
			// pass!  it's a reference a module import; no recursion to do.
		case false:
			if err := orderSteps_visit(wire.StepName, todo, loopDetector, result, m); err != nil {
				return nil
			}
		}
	}
	// Done: put this node in the results, and remove from todo set.
	*result = append(*result, node)
	delete(todo, node)
	return nil
}

type slotRefs []SlotReference

func (s slotRefs) Len() int      { return len(s) }
func (s slotRefs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s slotRefs) Less(i, j int) bool {
	return s[i].StepName < s[j].StepName || s[i].SlotName < s[j].SlotName
}

func inputSlotReferences(s StepUnion) (r []SlotReference) {
	switch x := s.(type) {
	case Operation:
		for k := range x.Inputs {
			r = append(r, k)
		}
		return r
	case Module:
		for _, imp := range x.Imports {
			switch imp2 := imp.(type) {
			case ImportRef_Catalog:
				// pass
			case ImportRef_Parent:
				r = append(r, SlotReference(imp2))
			case ImportRef_Ingest:
				// this is panic-worthy because it should've been checked earlier.
				panic("submodules can't have ingest imports!")
			}
		}
		return r
	}
	panic("unreachable")
}

func outputSlotReferences(s StepUnion) map[SlotName]struct{} {
	r := make(map[SlotName]struct{})
	switch x := s.(type) {
	case Operation:
		for k := range x.Outputs {
			r[k] = struct{}{}
		}
		return r
	case Module:
		for itemName := range x.Exports {
			r[SlotName(itemName)] = struct{}{}
		}
		return r
	}
	panic("unreachable")
}
