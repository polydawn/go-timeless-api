package funcs

import (
	"fmt"
	"sort"

	"go.polydawn.net/go-timeless-api"
)

type StepList []api.StepName
type StepTree []api.SubmoduleStepRef

/*
	AppendSubtree returns a new subtree, appending the given `t2` to this `t`
	after first contextualizing all the stepnames in `t2` with the given StepName.
	AppendSubtree is mostly used in ModuleOrderStepsDeep.
*/
func (t StepTree) AppendSubtree(submoduleName api.StepName, t2 StepTree) StepTree {
	l := len(t)
	t3 := make(StepTree, l+len(t2))
	copy(t3, t)
	for i, subStep := range t2 {
		t3[l+i] = subStep.Contextualize(api.SubmoduleRef(submoduleName))
	}
	return t3
}

/*
	DetachSubtree is roughly the opposite of AppendSubtree: it returns a new StepTree
	containing only those elements which were prefixed by `submoduleName`, and
	de-contextualizes them of that `submoduleName` prefix.
	This is often useful when traversing a module recursively while following a StepTree.
*/
func (t StepTree) DetachSubtree(submoduleName api.StepName) StepTree {
	t2 := StepTree{}
	for _, step := range t {
		if step.SubmoduleRef.First() == submoduleName {
			t2 = append(t2, step.Decontextualize())
		}
	}
	return t2
}

/*
	ModuleOrderStepsDeep is like ModuleOrderSteps, but returns *all* steps, recursively
	including all submodules and their steps.
*/
func ModuleOrderStepsDeep(m api.Module) (r StepTree, _ error) {
	levelOrder, err := ModuleOrderSteps(m)
	if err != nil {
		return nil, err
	}
	for _, stepName := range levelOrder {
		switch x := m.Steps[stepName].(type) {
		case api.Operation:
			r = append(r, api.SubmoduleStepRef{"", stepName})
		case api.Module:
			r = append(r, api.SubmoduleStepRef{"", stepName})
			subOrder, err := ModuleOrderStepsDeep(x)
			if err != nil {
				return nil, err
			}
			r = r.AppendSubtree(stepName, subOrder)
		}
	}
	return r, nil
}

/*
	ModuleOrderSteps returns a slice of StepName from a Module's Operations
	in order of topological sort based on SlotName each op references as inputs.

	We break ties based on lexigraphical sort on the step names.
	We choose this simple tie-breaker rather than attempting any fancier
	logic based on e.g. downstream dependencies, etc, because ease of
	understanding and the simplicity of predicting the result of the sort
	is more important than cleverness; so is the regional stability of the
	sort in the face of changes in other parts of the graph.
*/
func ModuleOrderSteps(m api.Module) (StepList, error) {
	// Alloc result accumulator.
	result := make([]api.StepName, 0, len(m.Steps))
	// Initialize todo set; it shrinks as we go.
	todo := make(map[api.StepName]struct{}, len(m.Steps))
	for step := range m.Steps {
		todo[step] = struct{}{}
	}
	// Sort operations by their name (this is our tiebreaker, in advance).
	stepsOrdered := make(stepNames, 0, len(m.Steps))
	for step := range m.Steps {
		stepsOrdered = append(stepsOrdered, step)
	}
	sort.Sort(stepsOrdered)
	// For each step: visit.  (This will recurse, and no-op itself internally as approrpriate for visited nodes.)
	for _, step := range stepsOrdered {
		if err := orderSteps_visit(step, todo, map[api.StepName]struct{}{}, &result, m); err != nil {
			return nil, err
		}
	}
	return result, nil
}

type stepNames []api.StepName

func (s stepNames) Len() int           { return len(s) }
func (s stepNames) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s stepNames) Less(i, j int) bool { return s[i] < s[j] }

// func OrderStepsYielding(m Module, want []SlotRef) ([]StepName, error)
//   Easier to implement this later by reverse walking the toposort of []StepName
//    that we've already computed and simply dropping irrelevant ones.
//    This approach has a lower CPU time if we've already done the toposort;
//    and we *always* want to have already done the toposort, because in final
//    UX we would always want our tool to return errors on invalid topology
//    even if when running in the "just this swatch of module" mode you
//    coincidentally wouldn't happen to hit the invalid region.

func orderSteps_visit(
	node api.StepName,
	todo map[api.StepName]struct{},
	loopDetector map[api.StepName]struct{},
	result *[]api.StepName,
	m api.Module,
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
	wires := inputSlotRefs(m.Steps[node])
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
			if op, ok := m.Steps[wire.StepName]; !ok {
				return fmt.Errorf("step %q has an invalid reference to %q: %q is not the name of a step in this module", node, wire, wire.StepName)
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
				return err
			}
		}
	}
	// Done: put this node in the results, and remove from todo set.
	*result = append(*result, node)
	delete(todo, node)
	return nil
}

type slotRefs []api.SlotRef

func (s slotRefs) Len() int      { return len(s) }
func (s slotRefs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s slotRefs) Less(i, j int) bool {
	return s[i].StepName < s[j].StepName || s[i].SlotName < s[j].SlotName
}

func inputSlotRefs(s api.StepUnion) (r []api.SlotRef) {
	switch x := s.(type) {
	case api.Operation:
		for _, ref := range x.Inputs {
			r = append(r, ref)
		}
		return r
	case api.Module:
		for _, imp := range x.Imports {
			switch imp2 := imp.(type) {
			case api.ImportRef_Catalog:
				// pass
			case api.ImportRef_Parent:
				r = append(r, api.SlotRef(imp2))
			case api.ImportRef_Ingest:
				// this is panic-worthy because it should've been checked earlier.
				panic("submodules can't have ingest imports!")
			}
		}
		return r
	}
	panic("unreachable")
}

func outputSlotReferences(s api.StepUnion) map[api.SlotName]struct{} {
	r := make(map[api.SlotName]struct{})
	switch x := s.(type) {
	case api.Operation:
		for k := range x.Outputs {
			r[k] = struct{}{}
		}
		return r
	case api.Module:
		for itemName := range x.Exports {
			r[api.SlotName(itemName)] = struct{}{}
		}
		return r
	}
	panic("unreachable")
}
