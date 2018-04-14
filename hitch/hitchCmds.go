/*
	Interfaces of hitch commands.

	The hitch CLI provides all these functions; it's also an interface used purely
	in-memory by some other systems (like heft) which handle both plan evaluation
	and release management in the same process.
*/
package hitch

import (
	"context"

	"go.polydawn.net/go-timeless-api"
)

type ViewCatalog func(
	context.Context,
	api.CatalogName,
) (*api.Catalog, error)
