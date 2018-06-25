package funcs

import (
	"context"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
	"go.polydawn.net/go-timeless-api/ingest"
)

type Pins map[api.SubmoduleSlotRef]api.WareID

func (t Pins) AppendSubtree(submoduleName api.StepName, t2 Pins) {
	for ref, wareID := range t2 {
		t[ref.Contextualize(api.SubmoduleRef(submoduleName))] = wareID
	}
}

func (t Pins) DetachSubtree(submoduleName api.StepName) Pins {
	t2 := Pins{}
	for ref, wareID := range t {
		if ref.First() == submoduleName {
			t2[ref.Decontextualize()] = wareID
		}
	}
	return t2
}

func ResolvePins(m api.Module, catalogTool hitch.ViewCatalogTool, ingestTool ingest.IngestTool) (Pins, error) {
	r := make(Pins)

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
			wareID, wareSourcing, err := ingestTool(context.TODO(), impRef2)
			if err != nil {
				return nil, err
			}
			r[api.SubmoduleSlotRef{"", api.SlotRef{"", slotName}}] = *wareID
			_ = wareSourcing // TODO
		}
	}

	// recurse, and fold all those references into our return set
	for stepName, step := range m.Steps {
		switch x := step.(type) {
		case api.Operation:
			// pass.  hakuna matata; operations only have local references to their module's imports.
		case api.Module:
			// recurse, and contextualize all refs from the deeper module(s).
			subPins, err := ResolvePins(x, catalogTool, nil)
			if err != nil {
				return nil, err
			}
			r.AppendSubtree(stepName, subPins)
		}
	}
	return r, nil
}
