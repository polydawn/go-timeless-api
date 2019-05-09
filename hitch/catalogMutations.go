package hitch

// This file is full of helpful mutation methods for module catalogs.

import (
	"fmt"

	"github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api"
)

// LineagePrependRelease returns a new modified lineage with the release pushed onto
// the top of the lineage's list of releases.
//
// Any error will be of category LookupError.
//
// A pointer is returned to express maybe-ness.
func LineagePrependRelease(lin api.Lineage, rel api.Release) (*api.Lineage, error) {
	// Check we're not about to insert a dupe name; reject if so.
	_, err := LineagePluckReleaseByName(lin, rel.Name)
	switch errcat.Category(err) {
	case nil:
		return nil, errcat.ErrorDetailed(ErrNameCollision,
			fmt.Sprintf("catalog %q already has a release named %q", lin.Name, rel.Name),
			map[string]string{
				"ref": api.ItemRef{lin.Name, rel.Name, ""}.String(),
			},
		)
	case ErrNoSuchRelease:
		// continue!
	default:
		return nil, err
	}
	// Allocate new array.  New stuff goes to top; old stuff to bottom.
	releases := make([]api.Release, len(lin.Releases)+1)
	copy(releases[1:], lin.Releases)
	releases[0] = rel
	lin.Releases = releases
	return &lin, nil
}
