package hitch

// This file is full of helpful mutation methods for module catalogs.

import (
	"fmt"

	"github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
)

// CatalogPrependRelease returns a new modified catalog with the release pushed onto
// the top of the catalog's list of releases.
//
// Any error will be of category LookupError.
//
// A pointer is returned to express maybe-ness.
func CatalogPrependRelease(modCat api.ModuleCatalog, rel api.Release) (*api.ModuleCatalog, error) {
	// Check we're not about to insert a dupe name; reject if so.
	_, err := CatalogPluckReleaseByName(modCat, rel.Name)
	switch errcat.Category(err) {
	case nil:
		return nil, errcat.ErrorDetailed(ErrNameCollision,
			fmt.Sprintf("catalog %q already has a release named %q", modCat.Name, rel.Name),
			map[string]string{
				"ref": api.ItemRef{modCat.Name, rel.Name, ""}.String(),
			},
		)
	case ErrNoSuchRelease:
		// continue!
	default:
		return nil, err
	}
	// Allocate new array.  New stuff goes to top; old stuff to bottom.
	releases := make([]api.Release, len(modCat.Releases)+1)
	copy(releases[1:], modCat.Releases)
	releases[0] = rel
	modCat.Releases = releases
	return &modCat, nil
}
