package api

// All 'parse_*' methods expect to be called on a struct that is all `ff_uninitialized` values at start, and mutate it.
// All 'string_*' methods return blanks if the value is `ff_uninitialized`.

import (
	"fmt"
	"strconv"
	"time"
)

func (ff *FilesetPackFilter) parse_uid(s string) (err error) {
	if ff.uid != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "keep":
		ff.uid = ff_keep
	default:
		ff.uid, err = strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("fileset filter: invalid option: uid must be 'keep' or an int")
		}
	}
	return nil
}
func (ff *FilesetPackFilter) parse_gid(s string) (err error) {
	if ff.gid != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "keep":
		ff.gid = ff_keep
	default:
		ff.gid, err = strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("fileset filter: invalid option: gid must be 'keep' or an int")
		}
	}
	return nil
}
func (ff *FilesetPackFilter) parse_mtime(s string) (err error) {
	if ff.mtime != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "keep":
		ff.mtime = ff_keep
	default:
		goto handle
	error:
		return fmt.Errorf("fileset filter: invalid option: mtime must be either 'keep', a unix timestamp integer beginning with '@', or an RFC3339 date string")
	handle:
		if len(s) == 0 {
			goto error
		}
		if s[0] == '@' {
			ff.mtime, err = strconv.ParseInt(s[1:], 10, 0)
			if err != nil {
				goto error
			}
			break
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			goto error
		}
		ff.mtime = t.Unix()
	}
	return nil
}
func (ff *FilesetPackFilter) parse_sticky(s string) (err error) {
	if ff.sticky != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "keep":
		ff.sticky = ff_keep
	case "ignore":
		ff.sticky = ff_ignore
	default:
		return fmt.Errorf("fileset filter: invalid option: sticky must be 'keep' or 'ignore'")
	}
	return nil
}
func (ff *FilesetPackFilter) parse_setid(s string) (err error) {
	if ff.setid != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "keep":
		ff.setid = ff_keep
	case "ignore":
		ff.setid = ff_ignore
	case "reject":
		ff.setid = ff_reject
	default:
		return fmt.Errorf("fileset filter: invalid option: setid must be 'keep', 'ignore', or 'reject'")
	}
	return nil
}
func (ff *FilesetPackFilter) parse_dev(s string) (err error) {
	if ff.dev != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "keep":
		ff.dev = ff_keep
	case "ignore":
		ff.dev = ff_ignore
	case "reject":
		ff.dev = ff_reject
	default:
		return fmt.Errorf("fileset filter: invalid option: dev must be 'keep', 'ignore', or 'reject'")
	}
	return nil
}

func (ff *FilesetPackFilter) string_uid() string {
	switch {
	case ff.uid == ff_unspecified:
		return ""
	case ff.uid == ff_keep:
		return "keep"
	case ff.uid >= 0:
		return strconv.Itoa(ff.uid)
	default:
		panic("invalid")
	}
}
func (ff *FilesetPackFilter) string_gid() string {
	switch {
	case ff.gid == ff_unspecified:
		return ""
	case ff.gid == ff_keep:
		return "keep"
	case ff.gid >= 0:
		return strconv.Itoa(ff.gid)
	default:
		panic("invalid")
	}
}
func (ff *FilesetPackFilter) string_mtime() string {
	switch {
	case ff.mtime == ff_unspecified:
		return ""
	case ff.mtime == ff_keep:
		return "keep"
	case ff.mtime >= 0:
		return "@" + strconv.FormatInt(ff.mtime, 10)
	default:
		panic("invalid")
	}
}
func (ff *FilesetPackFilter) string_sticky() string {
	switch {
	case ff.sticky == ff_unspecified:
		return ""
	case ff.sticky == ff_keep:
		return "keep"
	case ff.sticky == ff_ignore:
		return "ignore"
	default:
		panic("invalid")
	}
}
func (ff *FilesetPackFilter) string_setid() string {
	switch {
	case ff.setid == ff_unspecified:
		return ""
	case ff.setid == ff_keep:
		return "keep"
	case ff.setid == ff_ignore:
		return "ignore"
	case ff.setid == ff_reject:
		return "reject"
	default:
		panic("invalid")
	}
}
func (ff *FilesetPackFilter) string_dev() string {
	switch {
	case ff.dev == ff_unspecified:
		return ""
	case ff.dev == ff_keep:
		return "keep"
	case ff.dev == ff_ignore:
		return "ignore"
	case ff.dev == ff_reject:
		return "reject"
	default:
		panic("invalid")
	}
}

// ----

func (ff *FilesetUnpackFilter) parse_uid(s string) (err error) {
	if ff.uid != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "follow":
		ff.uid = ff_follow
	case "mine":
		ff.uid = ff_context
	default:
		ff.uid, err = strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("fileset filter: invalid option: uid must be 'follow', 'mine', or an int")
		}
	}
	return nil
}
func (ff *FilesetUnpackFilter) parse_gid(s string) (err error) {
	if ff.gid != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "follow":
		ff.gid = ff_follow
	case "mine":
		ff.gid = ff_context
	default:
		ff.gid, err = strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("fileset filter: invalid option: gid must be 'follow', 'mine' or an int")
		}
	}
	return nil
}
func (ff *FilesetUnpackFilter) parse_mtime(s string) (err error) {
	if ff.mtime != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "follow":
		ff.mtime = ff_follow
	case "now":
		ff.mtime = ff_context
	default:
		goto handle
	error:
		return fmt.Errorf("fileset filter: invalid option: mtime must be either 'follow', 'now', a unix timestamp integer beginning with '@', or an RFC3339 date string")
	handle:
		if len(s) == 0 {
			goto error
		}
		if s[0] == '@' {
			ff.mtime, err = strconv.ParseInt(s[1:], 10, 0)
			if err != nil {
				goto error
			}
			break
		}
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			goto error
		}
		ff.mtime = t.Unix()
	}
	return nil
}
func (ff *FilesetUnpackFilter) parse_sticky(s string) (err error) {
	if ff.sticky != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "follow":
		ff.sticky = ff_follow
	case "ignore":
		ff.sticky = ff_ignore
	default:
		return fmt.Errorf("fileset filter: invalid option: sticky must be 'follow' or 'ignore'")
	}
	return nil
}
func (ff *FilesetUnpackFilter) parse_setid(s string) (err error) {
	if ff.setid != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "follow":
		ff.setid = ff_follow
	case "ignore":
		ff.setid = ff_ignore
	case "reject":
		ff.setid = ff_reject
	default:
		return fmt.Errorf("fileset filter: invalid option: setid must be 'follow', 'ignore', or 'reject'")
	}
	return nil
}
func (ff *FilesetUnpackFilter) parse_dev(s string) (err error) {
	if ff.dev != ff_unspecified {
		return fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
	}
	switch s {
	case "follow":
		ff.dev = ff_follow
	case "ignore":
		ff.dev = ff_ignore
	case "reject":
		ff.dev = ff_reject
	default:
		return fmt.Errorf("fileset filter: invalid option: dev must be 'follow', 'ignore', or 'reject'")
	}
	return nil
}

func (ff *FilesetUnpackFilter) string_uid() string {
	switch {
	case ff.uid == ff_unspecified:
		return ""
	case ff.uid == ff_follow:
		return "follow"
	case ff.uid == ff_context:
		return "mine"
	case ff.uid >= 0:
		return strconv.Itoa(ff.uid)
	default:
		panic("invalid")
	}
}
func (ff *FilesetUnpackFilter) string_gid() string {
	switch {
	case ff.gid == ff_unspecified:
		return ""
	case ff.gid == ff_follow:
		return "follow"
	case ff.gid == ff_context:
		return "mine"
	case ff.gid >= 0:
		return strconv.Itoa(ff.gid)
	default:
		panic("invalid")
	}
}
func (ff *FilesetUnpackFilter) string_mtime() string {
	switch {
	case ff.mtime == ff_unspecified:
		return ""
	case ff.mtime == ff_follow:
		return "follow"
	case ff.mtime == ff_context:
		return "now"
	case ff.mtime >= 0:
		return "@" + strconv.FormatInt(ff.mtime, 10)
	default:
		panic("invalid")
	}
}
func (ff *FilesetUnpackFilter) string_sticky() string {
	switch {
	case ff.sticky == ff_unspecified:
		return ""
	case ff.sticky == ff_follow:
		return "follow"
	case ff.sticky == ff_ignore:
		return "ignore"
	default:
		panic("invalid")
	}
}
func (ff *FilesetUnpackFilter) string_setid() string {
	switch {
	case ff.setid == ff_unspecified:
		return ""
	case ff.setid == ff_follow:
		return "follow"
	case ff.setid == ff_ignore:
		return "ignore"
	case ff.setid == ff_reject:
		return "reject"
	default:
		panic("invalid")
	}
}
func (ff *FilesetUnpackFilter) string_dev() string {
	switch {
	case ff.dev == ff_unspecified:
		return ""
	case ff.dev == ff_follow:
		return "follow"
	case ff.dev == ff_ignore:
		return "ignore"
	case ff.dev == ff_reject:
		return "reject"
	default:
		panic("invalid")
	}
}
