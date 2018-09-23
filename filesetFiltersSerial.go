package api

import (
	"fmt"
	"strings"

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
			if err := ff.parse_uid(hunks[1]); err != nil {
				return ff, err
			}
		case "gid":
			if err := ff.parse_gid(hunks[1]); err != nil {
				return ff, err
			}
		case "mtime":
			if err := ff.parse_mtime(hunks[1]); err != nil {
				return ff, err
			}
		case "sticky":
			if err := ff.parse_sticky(hunks[1]); err != nil {
				return ff, err
			}
		case "setid":
			if err := ff.parse_setid(hunks[1]); err != nil {
				return ff, err
			}
		case "dev":
			if err := ff.parse_dev(hunks[1]); err != nil {
				return ff, err
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
	if s := x.string_uid(); s != "" {
		hunks = append(hunks, "uid="+s)
	}
	if s := x.string_gid(); s != "" {
		hunks = append(hunks, "gid="+s)
	}
	if s := x.string_mtime(); s != "" {
		hunks = append(hunks, "mtime="+s)
	}
	if s := x.string_sticky(); s != "" {
		hunks = append(hunks, "sticky="+s)
	}
	if s := x.string_setid(); s != "" {
		hunks = append(hunks, "setid="+s)
	}
	if s := x.string_dev(); s != "" {
		hunks = append(hunks, "dev="+s)
	}
	return strings.Join(hunks, ",")
}

var FilesetPackFilter_AsString_AtlasEntry = atlas.BuildEntry(FilesetPackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetPackFilter) (string, error) {
			return x.String(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseFilesetPackFilter)).
	Complete()

var FilesetPackFilter_AtlasEntry = atlas.BuildEntry(FilesetPackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetPackFilter) (map[string]string, error) {
			if x.initialized == false {
				return nil, nil
			}
			ffs := make(map[string]string, 6)
			if s := x.string_uid(); s != "" {
				ffs["uid"] = s
			}
			if s := x.string_gid(); s != "" {
				ffs["gid"] = s
			}
			if s := x.string_mtime(); s != "" {
				ffs["mtime"] = s
			}
			if s := x.string_sticky(); s != "" {
				ffs["sticky"] = s
			}
			if s := x.string_setid(); s != "" {
				ffs["setid"] = s
			}
			if s := x.string_dev(); s != "" {
				ffs["dev"] = s
			}
			return ffs, nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(x map[string]string) (FilesetPackFilter, error) {
			ff := FilesetPackFilter{true,
				ff_unspecified, ff_unspecified, ff_unspecified,
				ff_unspecified, ff_unspecified, ff_unspecified,
			}
			if x == nil {
				return ff, nil
			}
			if s, exists := x["uid"]; exists {
				if err := ff.parse_uid(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["gid"]; exists {
				if err := ff.parse_gid(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["mtime"]; exists {
				if err := ff.parse_mtime(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["sticky"]; exists {
				if err := ff.parse_sticky(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["setid"]; exists {
				if err := ff.parse_setid(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["dev"]; exists {
				if err := ff.parse_dev(s); err != nil {
					return ff, err
				}
			}
			return ff, nil
		})).
	Complete()

// ----

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
			if err := ff.parse_uid(hunks[1]); err != nil {
				return ff, err
			}
		case "gid":
			if err := ff.parse_gid(hunks[1]); err != nil {
				return ff, err
			}
		case "mtime":
			if err := ff.parse_mtime(hunks[1]); err != nil {
				return ff, err
			}
		case "sticky":
			if err := ff.parse_sticky(hunks[1]); err != nil {
				return ff, err
			}
		case "setid":
			if err := ff.parse_setid(hunks[1]); err != nil {
				return ff, err
			}
		case "dev":
			if err := ff.parse_dev(hunks[1]); err != nil {
				return ff, err
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
	if s := x.string_uid(); s != "" {
		hunks = append(hunks, "uid="+s)
	}
	if s := x.string_gid(); s != "" {
		hunks = append(hunks, "gid="+s)
	}
	if s := x.string_mtime(); s != "" {
		hunks = append(hunks, "mtime="+s)
	}
	if s := x.string_sticky(); s != "" {
		hunks = append(hunks, "sticky="+s)
	}
	if s := x.string_setid(); s != "" {
		hunks = append(hunks, "setid="+s)
	}
	if s := x.string_dev(); s != "" {
		hunks = append(hunks, "dev="+s)
	}
	return strings.Join(hunks, ",")
}

var FilesetUnpackFilter_AsString_AtlasEntry = atlas.BuildEntry(FilesetUnpackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetUnpackFilter) (string, error) {
			return x.String(), nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseFilesetUnpackFilter)).
	Complete()

var FilesetUnpackFilter_AtlasEntry = atlas.BuildEntry(FilesetUnpackFilter{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(
		func(x FilesetUnpackFilter) (map[string]string, error) {
			if x.initialized == false {
				return nil, nil
			}
			ffs := make(map[string]string, 6)
			if s := x.string_uid(); s != "" {
				ffs["uid"] = s
			}
			if s := x.string_gid(); s != "" {
				ffs["gid"] = s
			}
			if s := x.string_mtime(); s != "" {
				ffs["mtime"] = s
			}
			if s := x.string_sticky(); s != "" {
				ffs["sticky"] = s
			}
			if s := x.string_setid(); s != "" {
				ffs["setid"] = s
			}
			if s := x.string_dev(); s != "" {
				ffs["dev"] = s
			}
			return ffs, nil
		})).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(
		func(x map[string]string) (FilesetUnpackFilter, error) {
			ff := FilesetUnpackFilter{true,
				ff_unspecified, ff_unspecified, ff_unspecified,
				ff_unspecified, ff_unspecified, ff_unspecified,
			}
			if x == nil {
				return ff, nil
			}
			if s, exists := x["uid"]; exists {
				if err := ff.parse_uid(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["gid"]; exists {
				if err := ff.parse_gid(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["mtime"]; exists {
				if err := ff.parse_mtime(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["sticky"]; exists {
				if err := ff.parse_sticky(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["setid"]; exists {
				if err := ff.parse_setid(s); err != nil {
					return ff, err
				}
			}
			if s, exists := x["dev"]; exists {
				if err := ff.parse_dev(s); err != nil {
					return ff, err
				}
			}
			return ff, nil
		})).
	Complete()
