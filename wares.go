package api

/*
	Ware IDs are content-addressable, cryptographic hashes which uniquely identify
	a "ware" -- a packed filesystem snapshot.
	A ware contains one or more files and directories, and metadata for each.

	Ware IDs are serialized as a string in two parts, separated by a colon --
	for example like "git:f23ae1829" or "tar:WJL8or32vD".
	The first part communicates which kind of packing system computed the hash,
	and the second part is the hash itself.
*/
type WareID struct {
	Type PackType
	Hash string
}

/*
	A PackType string identifies what kind of packing format is used when
	packing a ware.  It's the first part of a WareID tuple.

	Typically, the desired PackType is an argument when using packing tools;
	whereas the PackType is communicated by the WareID when using unpack tools.

	PackTypes are a simple [a-zA-Z0-9] string.  Colons in particular are not
	allowable (since a PackType string is the first part of a WareID).
*/
type PackType string

// n.b. there is no typedef for the WareID.Hash, because
// *you never communicate them* outside of *the WareID* tuple:
// they're technically meaningless without having
// the PackType in hand to define their scope/encoding.
