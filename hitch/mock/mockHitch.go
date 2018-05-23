package mockhitch

import (
	"context"

	"github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

var (
	_ hitch.ViewCatalogTool = Fixture{}.ViewCatalog
)

type Fixture struct {
	Catalog map[api.ModuleName]api.ModuleCatalog
}

func (fix Fixture) ViewCatalog(
	_ context.Context,
	modName api.ModuleName,
) (*api.ModuleCatalog, error) {
	mcat, exists := fix.Catalog[modName]
	if !exists {
		return nil, errcat.Errorf(hitch.ErrModuleCatalogNotExist, "no catalog for module %q", modName)
	}
	return &mcat, nil
}
