package hitch

// This file is full of helpful traversal methods.
// It's considered a practical truism that catalogs are the organizational unit
// for actual on-disk storage, so generally those are loaded, and then we do
// traversals in-memory from there.

import (
	"fmt"

	"github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
)

// CatalogPluckReleaseByName traverses a module catalog and returns a release by name.
// This is a useful function because the serial form of release catalogs
// stores them in an ordered array.
//
// Any error will be of category LookupError.
func CatalogPluckReleaseByName(cat api.ModuleCatalog, releaseName api.ReleaseName) (*api.Release, error) {
	for _, rel := range cat.Releases {
		if rel.Name == releaseName {
			return &rel, nil
		}
	}
	return nil, errcat.ErrorDetailed(ErrNoSuchRelease,
		fmt.Sprintf("no such release %q in catalog %q", releaseName, cat.Name),
		map[string]string{
			"ref": api.ItemRef{cat.Name, releaseName, ""}.String(),
		},
	)
}

// CatalogPluckReleaseItem traverses a module catalog and returns a wareID looked up
// by release and item label.
//
// Any error will be of category LookupError.
func CatalogPluckReleaseItem(cat api.ModuleCatalog, releaseName api.ReleaseName, itemName api.ItemName) (*api.WareID, error) {
	rel, err := CatalogPluckReleaseByName(cat, releaseName)
	if err != nil {
		return nil, err
	}
	if item, ok := rel.Items[itemName]; ok {
		return &item, nil
	}
	return nil, errcat.ErrorDetailed(ErrNoSuchItem,
		fmt.Sprintf("no such item %q in release %q", itemName, api.ItemRef{cat.Name, releaseName, ""}.String()),
		map[string]string{
			"ref": api.ItemRef{cat.Name, releaseName, itemName}.String(),
		},
	)
}
