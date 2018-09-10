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
func ParseFilesetPackFilter(s string) (ff FilesetPackFilter, err error) {
	alreadyEncountered := make(map[string]bool, 6)
	for _, s := range strings.Split(s, ",") {
		hunks := strings.SplitN(strings.TrimSpace(s), "=", 2)
		if len(hunks) != 2 {
			return ff, fmt.Errorf("fileset filter: invalid format: must be a comma-separated list of k=v pairs")
		}
		if alreadyEncountered[hunks[0]] == true {
			return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
		}
		alreadyEncountered[hunks[0]] = true
		switch hunks[0] {
		case "uid":
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
			switch hunks[1] {
			case "keep":
				ff.sticky = ff_keep
			case "ignore":
				ff.sticky = ff_ignore
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: sticky must be 'keep' or 'ignore'")
			}
		case "setid":
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
	var sb strings.Builder
	sb.WriteString("uid=")
	switch {
	case x.uid == ff_keep:
		sb.WriteString("keep")
	case x.uid >= 0:
		sb.WriteString(strconv.Itoa(x.uid))
	default:
		panic("invalid")
	}
	sb.WriteString(",gid=")
	switch {
	case x.gid == ff_keep:
		sb.WriteString("keep")
	case x.gid >= 0:
		sb.WriteString(strconv.Itoa(x.gid))
	default:
		panic("invalid")
	}
	sb.WriteString(",mtime=")
	switch {
	case x.mtime == ff_keep:
		sb.WriteString("keep")
	case x.mtime >= 0:
		sb.WriteByte('@')
		sb.WriteString(strconv.FormatInt(x.mtime, 10))
	default:
		panic("invalid")
	}
	sb.WriteString(",sticky=")
	switch {
	case x.sticky == ff_keep:
		sb.WriteString("keep")
	case x.sticky == ff_ignore:
		sb.WriteString("ignore")
	default:
		panic("invalid")
	}
	sb.WriteString(",setid=")
	switch {
	case x.setid == ff_keep:
		sb.WriteString("keep")
	case x.setid == ff_ignore:
		sb.WriteString("ignore")
	case x.setid == ff_reject:
		sb.WriteString("reject")
	default:
		panic("invalid")
	}
	sb.WriteString(",dev=")
	switch {
	case x.dev == ff_keep:
		sb.WriteString("keep")
	case x.dev == ff_ignore:
		sb.WriteString("ignore")
	case x.dev == ff_reject:
		sb.WriteString("reject")
	default:
		panic("invalid")
	}
	return sb.String()
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
func ParseFilesetUnpackFilter(s string) (ff FilesetUnpackFilter, err error) {
	alreadyEncountered := make(map[string]bool, 6)
	for _, s := range strings.Split(s, ",") {
		hunks := strings.SplitN(strings.TrimSpace(s), "=", 2)
		if len(hunks) != 2 {
			return ff, fmt.Errorf("fileset filter: invalid format: must be a comma-separated list of k=v pairs")
		}
		if alreadyEncountered[hunks[0]] == true {
			return ff, fmt.Errorf("fileset filter: cannot specify the same option repeatedly")
		}
		alreadyEncountered[hunks[0]] = true
		switch hunks[0] {
		case "uid":
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
			switch hunks[1] {
			case "follow":
				ff.sticky = ff_follow
			case "ignore":
				ff.sticky = ff_ignore
			default:
				return ff, fmt.Errorf("fileset filter: invalid option: sticky must be 'follow' or 'ignore'")
			}
		case "setid":
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
	var sb strings.Builder
	sb.WriteString("uid=")
	switch {
	case x.uid == ff_follow:
		sb.WriteString("follow")
	case x.uid == ff_context:
		sb.WriteString("mine")
	case x.uid >= 0:
		sb.WriteString(strconv.Itoa(x.uid))
	default:
		panic("invalid")
	}
	sb.WriteString(",gid=")
	switch {
	case x.gid == ff_follow:
		sb.WriteString("follow")
	case x.gid == ff_context:
		sb.WriteString("mine")
	case x.gid >= 0:
		sb.WriteString(strconv.Itoa(x.gid))
	default:
		panic("invalid")
	}
	sb.WriteString(",mtime=")
	switch {
	case x.mtime == ff_follow:
		sb.WriteString("follow")
	case x.mtime == ff_context:
		sb.WriteString("now")
	case x.mtime >= 0:
		sb.WriteByte('@')
		sb.WriteString(strconv.FormatInt(x.mtime, 10))
	default:
		panic("invalid")
	}
	sb.WriteString(",sticky=")
	switch {
	case x.sticky == ff_follow:
		sb.WriteString("follow")
	case x.sticky == ff_ignore:
		sb.WriteString("ignore")
	default:
		panic("invalid")
	}
	sb.WriteString(",setid=")
	switch {
	case x.setid == ff_follow:
		sb.WriteString("follow")
	case x.setid == ff_ignore:
		sb.WriteString("ignore")
	case x.setid == ff_reject:
		sb.WriteString("reject")
	default:
		panic("invalid")
	}
	sb.WriteString(",dev=")
	switch {
	case x.dev == ff_follow:
		sb.WriteString("follow")
	case x.dev == ff_ignore:
		sb.WriteString("ignore")
	case x.dev == ff_reject:
		sb.WriteString("reject")
	default:
		panic("invalid")
	}
	return sb.String()
}

var FilesetUnpackFilter_AtlasEntry = atlas.BuildEntry(FilesetUnpackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetUnpackFilter) (string, error) {
			return x.String(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseFilesetUnpackFilter)).
	Complete()
