package api

/*
	WareID is a content-addressable, cryptographic hashes that uniquely identifies
	a "ware" -- a packed Fileset.
	(Fileset and Ware are distinct concepts because a fileset is not packed in
	any particular way and thus has no innate hash; a Ware is packed and hashed.)

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

// future: probably introduce a PackConfig type which is PackType+FilesetPackFilters.
// unclear what would be ergonomic for the unpack variation.  it's less used anyway.
