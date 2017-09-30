package api

/*
	This file is all serializable types used in Rio
	to define filesets, WareIDs, packing systems, and storage locations.
*/

import (
	"fmt"
	"strings"

	"github.com/polydawn/refmt/obj/atlas"
)

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

func ParseWareID(x string) (WareID, error) {
	if x == "" {
		return WareID{}, nil
	}
	ss := strings.SplitN(x, ":", 2)
	if len(ss) < 2 {
		return WareID{}, fmt.Errorf("wareIDs must have contain a colon character (they are of form \"<type>:<hash>\")")
	}
	return WareID{PackType(ss[0]), ss[1]}, nil
}

func (x WareID) String() string {
	switch {
	case x.Type == "":
		return ""
	case x.Hash == "":
		return string(x.Type) + ":-"
	default:
		return string(x.Type) + ":" + x.Hash
	}
}

var WareID_AtlasEntry = atlas.BuildEntry(WareID{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x WareID) (string, error) {
			return x.String(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(x string) (WareID, error) {
			return ParseWareID(x)
		})).
	Complete()

type AbsPath string // Identifier for output slots.  Coincidentally, a path.

type (
	/*
		WarehouseAddr strings describe a protocol and dial address for talking to
		a storage warehouse.

		The serial format is an opaque string, though they typically resemble
		(and for internal use, are parsed as) URLs.
	*/
	WarehouseAddr string

	/*
		Configuration details for a warehouse.

		Many warehouses don't *need* any configuration; the addr string
		can tell the whole story.  But if you need auth or other fanciness,
		here's the place to specify it.
	*/
	WarehouseCfg struct {
		Auth     string      // auth info, if needed.  usually points to another file.
		Addr     interface{} // additional addr info, for protocols that require it.
		Priority int         // higher is checked first.
	}
)

/*
	FilesetFilters define how certain filesystem metadata should be treated
	when packing or unpacking files.

	They are stored as strings for simplicity of API, but are more like enums:

	UID and GID can be one of:

		- blank -- meaning "default behavior" (differs for pack and unpack;
		   it's like "1000" for pack and "mine" for unpack).
		   (When Repeatr drives Rio, it defaults to "keep" for unpack.)
		- "keep" -- meaning preserve the attributes of the fileset (in packing)
		   or manifest exactly what the ware specifies (in unpacking).
		- "mine" -- which means to ignore the ware attributes and use the current
		   uid/gid instead (this is only valid for unpacking).
		- an integer -- which means to treat everything as having exactly that
		   numeric uid/gid.

	Mtime can be one of:

		- blank -- meaning "default behavior" (differs for pack and unpack;
		   it's like "@25000" for pack and "keep" for unpack).
		- "keep" -- meaning preserve the attributes of the fileset (in packing)
		   or manifest exactly what the ware specifies (in unpacking).
		- "@" + an integer -- which means to set all times to the integer,
		   interpreted as a unix timestamp.
		- an RFC3339 date -- which means to set all times to that date
		   (and note that this will *not* survive serialization as such again;
		   it will be converted to the "@unix" format).

	Sticky is a simple bool: if true, setuid/setgid/sticky bits will be preserved
	on unpack.  The sticky bool has no meaning on pack; those bits are always packed.
	Repeatr always sets the sticky bool to true when using Rio,
	but it defaults to false when using Rio's command line.
	(For comparison, your system tar command tends to do the same:
	sticky bits are not unpacked by default because of the security implications
	if the user is unwary.)

	Despite being a simple bool, Stick is still serialized as a string:
	either "keep" or "zero".  A blank string is used as a ternary value that
	means "default", contextually -- this is necessary for communicating
	filter settings through layers of the stack without ambiguity in deeper
	layers about whether or not a user actually requested a specific setting.
*/
type FilesetFilters struct {
	Uid    string `refmt:"uid,omitempty"`
	Gid    string `refmt:"gid,omitempty"`
	Mtime  string `refmt:"mtime,omitempty"`
	Sticky string `refmt:"sticky,omitempty"`
}

var FilesetFilters_AtlasEntry = atlas.BuildEntry(FilesetFilters{}).StructMap().Autogenerate().Complete()

var (
	Filter_NoMutation     = FilesetFilters{"keep", "keep", "keep", "keep"}                 // The default filters on repeatr inputs.
	Filter_DefaultFlatten = FilesetFilters{"1000", "1000", "2010-01-01T00:00:00Z", "keep"} // The default filters on repeatr outputs and rio pack.
	Filter_LowPriv        = FilesetFilters{"mine", "mine", "keep", "zero"}                 // The default filters on rio unpack.
)
