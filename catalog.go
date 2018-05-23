package api

type ModuleName string
type ReleaseName string
type ItemName string

// n.b. if you're a developer following the name changes from earlier TLapi:
//  we're pivoting to refer to "the catalog" as *all* the things from *all* the projects;
//  so the all the releases for one project are now a "module catalog"
//  (read as: "one module's catalog in the whole catalog").

// ModuleCatalog contains the metadata for all releases for a particular module.
// Treat it as an append-only record: new releases append to the module's catalog.
type ModuleCatalog struct {
	// Name of self.
	Name ModuleName

	// Ordered list of release entries.
	// Order not particularly important, though UIs generally display in this order.
	// Most recent entries are should be placed at the top (e.g. index zero).
	//
	// Each entry must have a unique ReleaseName in the scope of this ModuleCatalog.
	Releases []Release
}

// Release describes a single atomic release of wares.
// Each release must have a name, and contains a set of items, where each item
// refers to a WareID.
//
// Releases are used to group something chronologically; items in a release
// are used to distinguish between multiple artifacts in a release.
//
// In the context of building software, a Release usually has semantics lining
// up with "a bunch of software built from a particular source checkout".
// And thus, typically, there is also an Item in the release called "src";
// and often enough, this will be a "git" wareID.
// Other Item names likely to appear might be "linux-amd64", for example.
// All of this is convention, however; releases could just as well be used
// to track various versions of a photo album.
//
// It is recommended that a series of Release entries in a ModuleCatalog
// should stick to the same set of ItemName over time, because consumers
// of catalog information generally expect this, and changing Item names
// may produce work for other people.
type Release struct {
	Name     ReleaseName
	Items    map[ItemName]WareID
	Metadata map[string]string
	Hazards  map[string]string
}
