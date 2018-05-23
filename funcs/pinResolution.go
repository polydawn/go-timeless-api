package funcs

import (
	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

func ResolvePins(m api.Module, catalogTool hitch.ViewCatalogTool) (r map[api.SubmoduleSlotReference]api.WareID, _ error) {
	for stepName, step := range m.Operations {
		switch x := step.(type) {
		case api.Operation:
			// TODO
			//for ref, imp := range x.ImportsResolved {
			//	catalogTool(context.Background(), imp.???)
			//}
		case api.Module:
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
