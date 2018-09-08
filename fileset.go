package api

import (
	"time"
)

type FilesetPackFilter struct {
	uid    int   // keep, [int]
	gid    int   // keep, [int]
	mtime  int64 // keep, [value] // we *could* support a 'now' mode, but we're all about discouraging that kind of nonsense in the fileset.
	sticky int   // keep, ignore // i don't actually know why you'd ever want to zero out a sticky bit, but it's here for completeness.
	setid  int   // keep, ignore, reject
	dev    int   // keep, ignore, reject
}

type FilesetUnpackFilter struct {
	uid    int   // follow, mine, [int]
	gid    int   // follow, mine, [int]
	mtime  int64 // follow, now, [value]
	sticky int   // follow, ignore
	setid  int   // follow, ignore, reject
	dev    int   // follow, ignore, reject
}

var (
	FilesetPackFilter_Lossless     = FilesetPackFilter{ff_keep, ff_keep, ff_keep, ff_keep, ff_keep, ff_keep}   // The default filters on... nothing, really.
	FilesetPackFilter_Flatten      = FilesetPackFilter{1000, 1000, defaultTime, ff_keep, ff_keep, ff_keep}     // The default filters on repeatr outputs.
	FilesetPackFilter_Conservative = FilesetPackFilter{1000, 1000, defaultTime, ff_keep, ff_reject, ff_reject} // The default filters on rio pack.  Guides you away from anything that would require privs to unpack again.

	FilesetUnpackFilter_Lossless = FilesetUnpackFilter{ff_follow, ff_follow, ff_follow, ff_follow, ff_follow, ff_follow}   // The default filters on repeatr inputs.
	FilesetUnpackFilter_LowPriv  = FilesetUnpackFilter{ff_context, ff_context, ff_follow, ff_follow, ff_reject, ff_reject} // The default filters on rio unpack.

	// note that the 'ignore' modes are never used in any of our common defaults.  they're only there for the user realizes they want them and require opt in.
)

var defaultTime int64 = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

const (
	ff_keep    = -2
	ff_follow  = -2
	ff_ignore  = -3
	ff_reject  = -4 // if trying to figure out caching, can map this into "ignore".
	ff_context = -5 // if trying to figure out caching, must map this into a real value.
)
