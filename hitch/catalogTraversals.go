package hitch

// This file is full of helpful traversal methods.
// It's considered a practical truism that lineages are the organizational unit
// for actual on-disk storage, so generally those are loaded, and then we do
// traversals in-memory from there.

import (
	"fmt"

	"github.com/warpfork/go-errcat"

	api "github.com/polydawn/go-timeless-api"
)

// LineagePluckReleaseByName traverses a linage and returns a release by name.
// This is a useful function because the serial form of release catalogs
// stores them in an ordered array.
//
// An error may be returned of category LookupError.
//
// A pointer is returned to express maybe-ness; mutating it has no effect.
func LineagePluckReleaseByName(lin api.Lineage, releaseName api.ReleaseName) (*api.Release, error) {
	for _, rel := range lin.Releases {
		if rel.Name == releaseName {
			return &rel, nil
		}
	}
	return nil, errcat.ErrorDetailed(ErrNoSuchRelease,
		fmt.Sprintf("no such release %q in lineage %q", releaseName, lin.Name),
		map[string]string{
			"ref": api.ItemRef{lin.Name, releaseName, ""}.String(),
		},
	)
}

// LineagePluckReleaseItem traverses a lineage and returns a wareID looked up
// by release and item label.
//
// An error may be returned of category LookupError.
//
// A pointer is returned to express maybe-ness; mutating it has no effect.
func LineagePluckReleaseItem(lin api.Lineage, releaseName api.ReleaseName, itemName api.ItemName) (*api.WareID, error) {
	rel, err := LineagePluckReleaseByName(lin, releaseName)
	if err != nil {
		return nil, err
	}
	if item, ok := rel.Items[itemName]; ok {
		return &item, nil
	}
	return nil, errcat.ErrorDetailed(ErrNoSuchItem,
		fmt.Sprintf("no such item %q in release %q", itemName, api.ItemRef{lin.Name, releaseName, ""}.String()),
		map[string]string{
			"ref": api.ItemRef{lin.Name, releaseName, itemName}.String(),
		},
	)
}
