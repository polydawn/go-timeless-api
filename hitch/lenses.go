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

// FocusReleaseByName traverses a catalog and returns a release by name.
// This is a useful function because the serial form of release catalogs
// stores them in an ordered array.
//
// Any error will be of category LookupError.
func PluckReleaseByName(cat api.Catalog, releaseName api.ReleaseName) (*api.ReleaseEntry, error) {
	for _, rel := range cat.Releases {
		if rel.Name == releaseName {
			return &rel, nil
		}
	}
	return nil, errcat.ErrorDetailed(ErrNoSuchRelease,
		fmt.Sprintf("no such release %q in catalog %q", releaseName, cat.Name),
		map[string]string{
			"name": api.ReleaseItemID{cat.Name, releaseName, ""}.String(),
		},
	)
}

// PluckReleaseItem traverses a catalog and returns a wareID looked up by
// release and item label.
//
// Any error will be of category LookupError.
func PluckReleaseItem(cat api.Catalog, releaseName api.ReleaseName, itemName api.ItemName) (*api.WareID, error) {
	rel, err := PluckReleaseByName(cat, releaseName)
	if err != nil {
		return nil, err
	}
	if item, ok := rel.Items[itemName]; ok {
		return &item, nil
	}
	return nil, errcat.ErrorDetailed(ErrNoSuchRelease,
		fmt.Sprintf("no such item %q in release %q", itemName, api.ReleaseItemID{cat.Name, releaseName, ""}.String()),
		map[string]string{
			"name": api.ReleaseItemID{cat.Name, releaseName, itemName}.String(),
		},
	)
}

type LookupError string

const (
	ErrNoSuchCatalog LookupError = ("no-such-catalog")
	ErrNoSuchRelease LookupError = ("no-such-release")
	ErrNoSuchItem    LookupError = ("no-such-item")
)
