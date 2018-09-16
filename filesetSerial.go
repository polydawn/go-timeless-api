package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/polydawn/refmt/obj/atlas"
)

func MustParseFilesetPackFilter(s string) FilesetPackFilter {
	ff, err := ParseFilesetPackFilter(s)
	if err != nil {
		panic(err)
	}
	return ff
}
func ParseFilesetPackFilter(s string) (_ FilesetPackFilter, err error) {
	ff := FilesetPackFilter{true,
		ff_unspecified, ff_unspecified, ff_unspecified,
		ff_unspecified, ff_unspecified, ff_unspecified,
	}
	if s == "" {
		return ff, nil
	}
	for _, s := range strings.Split(s, ",") {
		hunks := strings.SplitN(strings.TrimSpace(s), "=", 2)
		if len(hunks) != 2 {
			return ff, fmt.Errorf("fileset filter: invalid format: must be a comma-separated list of k=v pairs")
		}
		switch hunks[0] {
		case "uid":
			if ff.uid != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "keep":
				ff.uid = ff_keep
			default:
				ff.uid, err = strconv.Atoi(hunks[1])
				if err != nil {
					return ff, fmt.Errorf("fileset filter: invalid option: uid must be 'keep' or an int")
				}
			}
		case "gid":
			if ff.gid != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "keep":
				ff.gid = ff_keep
			default:
				ff.gid, err = strconv.Atoi(hunks[1])
				if err != nil {
					return ff, fmt.Errorf("fileset filter: invalid option: gid must be 'keep' or an int")
				}
			}
		case "mtime":
			if ff.mtime != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "keep":
				ff.mtime = ff_keep
			default:
				goto handle
			error:
				return ff, fmt.Errorf("fileset filter: invalid option: mtime must be either 'keep', a unix timestamp integer beginning with '@', or an RFC3339 date string")
			handle:
				if len(hunks[1]) == 0 {
					goto error
				}
				if hunks[1][0] == '@' {
					ff.mtime, err = strconv.ParseInt(hunks[1][1:], 10, 0)
					if err != nil {
						goto error
					}
					break
				}
				t, err := time.Parse(time.RFC3339, hunks[1])
				if err != nil {
					goto error
				}
				ff.mtime = t.Unix()
			}
		case "sticky":
			if ff.sticky != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "keep":
				ff.sticky = ff_keep
			case "ignore":
				ff.sticky = ff_ignore
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: sticky must be 'keep' or 'ignore'")
			}
		case "setid":
			if ff.setid != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "keep":
				ff.setid = ff_keep
			case "ignore":
				ff.setid = ff_ignore
			case "reject":
				ff.setid = ff_reject
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: setid must be 'keep', 'ignore', or 'reject'")
			}
		case "dev":
			if ff.dev != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "keep":
				ff.dev = ff_keep
			case "ignore":
				ff.dev = ff_ignore
			case "reject":
				ff.dev = ff_reject
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: dev must be 'keep', 'ignore', or 'reject'")
			}
		default:
			return ff, fmt.Errorf("fileset filter: unknown option: %q is not recognized", hunks[0])
		}
	}
	return ff, nil
}

func (x FilesetPackFilter) String() (v string) {
	if x.initialized == false {
		return ""
	}
	hunks := make([]string, 0, 6)
	switch {
	case x.uid == ff_unspecified:
		// skip
	case x.uid == ff_keep:
		hunks = append(hunks, "uid=keep")
	case x.uid >= 0:
		hunks = append(hunks, "uid="+strconv.Itoa(x.uid))
	default:
		panic("invalid")
	}
	switch {
	case x.gid == ff_unspecified:
		// skip
	case x.gid == ff_keep:
		hunks = append(hunks, "gid=keep")
	case x.gid >= 0:
		hunks = append(hunks, "gid="+strconv.Itoa(x.gid))
	default:
		panic("invalid")
	}
	switch {
	case x.mtime == ff_unspecified:
		// skip
	case x.mtime == ff_keep:
		hunks = append(hunks, "mtime=keep")
	case x.mtime >= 0:
		hunks = append(hunks, "mtime=@"+strconv.FormatInt(x.mtime, 10))
	default:
		panic("invalid")
	}
	switch {
	case x.sticky == ff_unspecified:
		// skip
	case x.sticky == ff_keep:
		hunks = append(hunks, "sticky=keep")
	case x.sticky == ff_ignore:
		hunks = append(hunks, "sticky=ignore")
	default:
		panic("invalid")
	}
	switch {
	case x.setid == ff_unspecified:
		// skip
	case x.setid == ff_keep:
		hunks = append(hunks, "setid=keep")
	case x.setid == ff_ignore:
		hunks = append(hunks, "setid=ignore")
	case x.setid == ff_reject:
		hunks = append(hunks, "setid=reject")
	default:
		panic("invalid")
	}
	switch {
	case x.dev == ff_unspecified:
		// skip
	case x.dev == ff_keep:
		hunks = append(hunks, "dev=keep")
	case x.dev == ff_ignore:
		hunks = append(hunks, "dev=ignore")
	case x.dev == ff_reject:
		hunks = append(hunks, "dev=reject")
	default:
		panic("invalid")
	}
	return strings.Join(hunks, ",")
}

var FilesetPackFilter_AtlasEntry = atlas.BuildEntry(FilesetPackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetPackFilter) (string, error) {
			return x.String(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseFilesetPackFilter)).
	Complete()

func MustParseFilesetUnpackFilter(s string) FilesetUnpackFilter {
	ff, err := ParseFilesetUnpackFilter(s)
	if err != nil {
		panic(err)
	}
	return ff
}
func ParseFilesetUnpackFilter(s string) (_ FilesetUnpackFilter, err error) {
	ff := FilesetUnpackFilter{true,
		ff_unspecified, ff_unspecified, ff_unspecified,
		ff_unspecified, ff_unspecified, ff_unspecified,
	}
	if s == "" {
		return ff, nil
	}
	for _, s := range strings.Split(s, ",") {
		hunks := strings.SplitN(strings.TrimSpace(s), "=", 2)
		if len(hunks) != 2 {
			return ff, fmt.Errorf("fileset filter: invalid format: must be a comma-separated list of k=v pairs")
		}
		switch hunks[0] {
		case "uid":
			if ff.uid != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "follow":
				ff.uid = ff_follow
			case "mine":
				ff.uid = ff_context
			default:
				ff.uid, err = strconv.Atoi(hunks[1])
				if err != nil {
					return ff, fmt.Errorf("fileset filter: invalid option: uid must be 'follow', 'mine', or an int")
				}
			}
		case "gid":
			if ff.gid != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "follow":
				ff.gid = ff_follow
			case "mine":
				ff.gid = ff_context
			default:
				ff.gid, err = strconv.Atoi(hunks[1])
				if err != nil {
					return ff, fmt.Errorf("fileset filter: invalid option: gid must be 'follow', 'mine' or an int")
				}
			}
		case "mtime":
			if ff.mtime != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "follow":
				ff.mtime = ff_follow
			case "now":
				ff.mtime = ff_context
			default:
				goto handle
			error:
				return ff, fmt.Errorf("fileset filter: invalid option: mtime must be either 'follow', 'now', a unix timestamp integer beginning with '@', or an RFC3339 date string")
			handle:
				if len(hunks[1]) == 0 {
					goto error
				}
				if hunks[1][0] == '@' {
					ff.mtime, err = strconv.ParseInt(hunks[1][1:], 10, 0)
					if err != nil {
						goto error
					}
					break
				}
				t, err := time.Parse(time.RFC3339, hunks[1])
				if err != nil {
					goto error
				}
				ff.mtime = t.Unix()
			}
		case "sticky":
			if ff.sticky != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "follow":
				ff.sticky = ff_follow
			case "ignore":
				ff.sticky = ff_ignore
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: sticky must be 'follow' or 'ignore'")
			}
		case "setid":
			if ff.setid != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "follow":
				ff.setid = ff_follow
			case "ignore":
				ff.setid = ff_ignore
			case "reject":
				ff.setid = ff_reject
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: setid must be 'follow', 'ignore', or 'reject'")
			}
		case "dev":
			if ff.dev != ff_unspecified {
				return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
			}
			switch hunks[1] {
			case "follow":
				ff.dev = ff_follow
			case "ignore":
				ff.dev = ff_ignore
			case "reject":
				ff.dev = ff_reject
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: dev must be 'follow', 'ignore', or 'reject'")
			}
		default:
			return ff, fmt.Errorf("fileset filter: unknown option: %q is not recognized", hunks[0])
		}
	}
	return ff, nil
}

func (x FilesetUnpackFilter) String() (v string) {
	if x.initialized == false {
		return ""
	}
	hunks := make([]string, 0, 6)
	switch {
	case x.uid == ff_unspecified:
		// skip
	case x.uid == ff_follow:
		hunks = append(hunks, "uid=follow")
	case x.uid == ff_context:
		hunks = append(hunks, "uid=mine")
	case x.uid >= 0:
		hunks = append(hunks, "uid="+strconv.Itoa(x.uid))
	default:
		panic("invalid")
	}
	switch {
	case x.gid == ff_unspecified:
		// skip
	case x.gid == ff_follow:
		hunks = append(hunks, "gid=follow")
	case x.gid == ff_context:
		hunks = append(hunks, "gid=mine")
	case x.gid >= 0:
		hunks = append(hunks, "gid="+strconv.Itoa(x.gid))
	default:
		panic("invalid")
	}
	switch {
	case x.mtime == ff_unspecified:
		// skip
	case x.mtime == ff_follow:
		hunks = append(hunks, "mtime=follow")
	case x.mtime == ff_context:
		hunks = append(hunks, "mtime=now")
	case x.mtime >= 0:
		hunks = append(hunks, "mtime=@"+strconv.FormatInt(x.mtime, 10))
	default:
		panic("invalid")
	}
	switch {
	case x.sticky == ff_unspecified:
		// skip
	case x.sticky == ff_follow:
		hunks = append(hunks, "sticky=follow")
	case x.sticky == ff_ignore:
		hunks = append(hunks, "sticky=ignore")
	default:
		panic("invalid")
	}
	switch {
	case x.setid == ff_unspecified:
		// skip
	case x.setid == ff_follow:
		hunks = append(hunks, "setid=follow")
	case x.setid == ff_ignore:
		hunks = append(hunks, "setid=ignore")
	case x.setid == ff_reject:
		hunks = append(hunks, "setid=reject")
	default:
		panic("invalid")
	}
	switch {
	case x.dev == ff_unspecified:
		// skip
	case x.dev == ff_follow:
		hunks = append(hunks, "dev=follow")
	case x.dev == ff_ignore:
		hunks = append(hunks, "dev=ignore")
	case x.dev == ff_reject:
		hunks = append(hunks, "dev=reject")
	default:
		panic("invalid")
	}
	return strings.Join(hunks, ",")
}

var FilesetUnpackFilter_AtlasEntry = atlas.BuildEntry(FilesetUnpackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetUnpackFilter) (string, error) {
			return x.String(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseFilesetUnpackFilter)).
	Complete()
