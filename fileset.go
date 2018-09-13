package api

import (
	"time"
)

// Note that these properties correlate closely with what you'll see in
//  Rio's FS packages: go.polydawn.net/rio/fs.Metadata in particular
//   is a struct describing the same properties that these filters modify.

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
	FilesetPackFilter_Flatten      = FilesetPackFilter{1000, 1000, DefaultTime, ff_keep, ff_keep, ff_keep}     // The default filters on repeatr outputs.
	FilesetPackFilter_Conservative = FilesetPackFilter{1000, 1000, DefaultTime, ff_keep, ff_reject, ff_reject} // The default filters on rio pack.  Guides you away from anything that would require privs to unpack again.

	FilesetUnpackFilter_Lossless = FilesetUnpackFilter{ff_follow, ff_follow, ff_follow, ff_follow, ff_follow, ff_follow}   // The default filters on repeatr inputs.
	FilesetUnpackFilter_LowPriv  = FilesetUnpackFilter{ff_context, ff_context, ff_follow, ff_follow, ff_reject, ff_reject} // The default filters on rio unpack.

	// note that the 'ignore' modes are never used in any of our common defaults.  they're only there for the user realizes they want them and require opt in.
)

var DefaultTime int64 = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

const (
	ff_unspecified = -1 // means not configured.  serialize as such; cannot use, must stack with defaults first.
	ff_keep        = -2
	ff_follow      = -2
	ff_ignore      = -3
	ff_reject      = -4 // if trying to figure out caching, can map this into "ignore".
	ff_context     = -5 // if trying to figure out caching, must map this into a real value.
)

func (ff FilesetPackFilter) Uid() (keep bool, setTo int) {
	return ff.uid == ff_keep, ff.uid
}
func (ff FilesetPackFilter) Gid() (keep bool, setTo int) {
	return ff.gid == ff_keep, ff.gid
}
func (ff FilesetPackFilter) Mtime() (keep bool, setTo time.Time) {
	return ff.mtime == ff_keep, time.Unix(ff.mtime, 0)
}
func (ff FilesetPackFilter) MtimeUnix() (keep bool, setTo int64) {
	return ff.mtime == ff_keep, ff.mtime
}
func (ff FilesetPackFilter) Sticky() (keep bool) {
	return ff.sticky == ff_keep
}
func (ff FilesetPackFilter) Setid() (keep bool, reject bool) {
	return ff.setid == ff_keep, ff.setid == ff_reject
}
func (ff FilesetPackFilter) Dev() (keep bool, reject bool) {
	return ff.dev == ff_keep, ff.dev == ff_reject
}

func (ff FilesetUnpackFilter) Uid() (follow, setMine bool, setTo int) {
	return ff.uid == ff_follow, ff.uid == ff_context, ff.uid
}
func (ff FilesetUnpackFilter) Gid() (follow, setMine bool, setTo int) {
	return ff.gid == ff_follow, ff.gid == ff_context, ff.gid
}
func (ff FilesetUnpackFilter) Mtime() (follow, setNow bool, setTo time.Time) {
	return ff.mtime == ff_follow, ff.mtime == ff_context, time.Unix(ff.mtime, 0)
}
func (ff FilesetUnpackFilter) MtimeUnix() (follow, now bool, setTo int64) {
	return ff.mtime == ff_follow, ff.mtime == ff_context, ff.mtime
}
func (ff FilesetUnpackFilter) Sticky() (follow bool) {
	return ff.sticky == ff_follow
}
func (ff FilesetUnpackFilter) Setid() (follow bool, reject bool) {
	return ff.setid == ff_follow, ff.setid == ff_reject
}
func (ff FilesetUnpackFilter) Dev() (follow bool, reject bool) {
	return ff.dev == ff_follow, ff.dev == ff_reject
}
func (ff FilesetUnpackFilter) Altering() bool {
	return ff.uid != ff_follow ||
		ff.gid != ff_follow ||
		ff.mtime != ff_follow ||
		ff.sticky != ff_follow ||
		(ff.setid != ff_follow && ff.setid != ff_reject) ||
		(ff.dev != ff_follow && ff.dev != ff_reject)
}
