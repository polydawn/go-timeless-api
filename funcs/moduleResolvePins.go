package funcs

import (
	"context"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

func ResolvePins(m api.Module, catalogTool hitch.ViewCatalogTool) (map[api.SubmoduleSlotRef]api.WareID, error) {
	r := make(map[api.SubmoduleSlotRef]api.WareID)

	// resolve each of our imports in this module
	for slotName, impRef := range m.Imports {
		_, _ = slotName, impRef
		switch impRef2 := impRef.(type) {
		case api.ImportRef_Catalog:
			mcat, err := catalogTool(context.TODO(), impRef2.ModuleName)
			if err != nil {
				return nil, err
			}
			wareID, err := CatalogPluckReleaseItem(*mcat, impRef2.ReleaseName, impRef2.ItemName)
			if err != nil {
				return nil, err
			}
			r[api.SubmoduleSlotRef{"", api.SlotRef{"", slotName}}] = *wareID
		case api.ImportRef_Parent:
			// pass.  we don't resolve these in advance; and it's checked by the 'OrderSteps' func that this refers to *something*.
		case api.ImportRef_Ingest:
			// TODO we need an 'ingestTool' of some kind, at least for the topmost module.
		}
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
