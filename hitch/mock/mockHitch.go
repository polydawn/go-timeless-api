package mockhitch

import (
	"context"

	"github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/hitch"
)

var (
	_ hitch.ViewLineageTool = Fixture{}.ViewLineage
)

type Fixture struct {
	Catalog map[api.ModuleName]api.Lineage
}

func (fix Fixture) ViewLineage(
	_ context.Context,
	modName api.ModuleName,
) (*api.Lineage, error) {
	mcat, exists := fix.Catalog[modName]
	if !exists {
		return nil, errcat.Errorf(hitch.ErrNoSuchLineage, "no lineage for module %q", modName)
	}
	return &mcat, nil
}
