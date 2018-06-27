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

type ViewCatalogTool func(
	context.Context,
	api.ModuleName,
) (*api.ModuleCatalog, error)

type ViewWarehousesTool func(
	context.Context,
	api.ModuleName,
) (*api.WareSourcing, error)

// for example of some other likely funcs coming up:

//type ViewReplayTool func(
//	context.Context,
//	api.ModuleName,
//	api.ReleaseName,
//) (*api.ModuleReplay, error)

// n.b. this doens't really give full control over insertion order
// but that's not something we really expect to do over an api;
// hitch cli does have such powers, though.
//type WriteReleaseTool func(
//	context.Context,
//	api.ModuleName,
//	api.ReleaseName,
//	map[ItemName]WareID,
//	optionallyReplay *api.ModuleReplay,
//	optionallySourcingHints *api.WareSourcing,
//)
