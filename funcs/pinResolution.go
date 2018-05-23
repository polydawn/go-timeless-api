package funcs

import (
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

func ResolvePins(m api.Module, catalogTool hitch.ViewCatalogTool) (map[api.SubmoduleSlotReference]api.WareID, error) {
	r := make(map[api.SubmoduleSlotReference]api.WareID)

	// resolve each of our imports in this module
	for slotName, impRef := range m.Imports {
		// TODO this'll depend on the type of import; let's make that type better defined first.
		// TODO and we also might need an 'ingestTool' of some kind, at least for the topmost module.
		_, _ = slotName, impRef
	}

	// recurse, and fold all those references into our return set
	for stepName, step := range m.Operations {
		switch x := step.(type) {
		case api.Operation:
			// pass.  hakuna matata; operations only have local references to their module's imports.
		case api.Module:
			// recurse, and contextualize all refs from the deeper module(s).
			subPins, err := ResolvePins(x, catalogTool)
			if err != nil {
				return nil, err
			}
			for ref, wareID := range subPins {
				r[ref.Contextualize(stepName)] = wareID
			}
		}
	}
	return r, nil
}
