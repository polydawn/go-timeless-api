package apiutil

import (
	"fmt"
	"strconv"
	"time"

	"go.polydawn.net/go-timeless-api"
)

/*
	Merges 'defaults' into 'user' filters, setting any blanks in the 'user'
	filters to the value from 'defaults', and returning the result.
*/
func MergeFilters(user api.FilesetFilters, defaults api.FilesetFilters) api.FilesetFilters {
	if user.Uid == "" {
		user.Uid = defaults.Uid
	}
	if user.Gid == "" {
		user.Gid = defaults.Gid
	}
	if user.Mtime == "" {
		user.Mtime = defaults.Mtime
	}
	if user.Sticky == "" {
		user.Sticky = defaults.Sticky
	}
	return user
}

type FilterPurpose bool

const (
	FilterPurposePack   FilterPurpose = false
	FilterPurposeUnpack FilterPurpose = true
)

const (
	FilterKeep = -1
	FilterMine = -2
)

/*
	Process the serializable API FilesetFilters into a more easily-used format,
	and return any errors with validating the API strings.
*/
func ProcessFilters(ff api.FilesetFilters, mode FilterPurpose) (uf FilesetFilters, err error) {
	// Parse UID.
	switch ff.Uid {
	case "":
		switch mode {
		case FilterPurposePack:
			uf.Uid = DefaultUid
		case FilterPurposeUnpack:
			uf.Uid = FilterMine
		}
	case "keep":
		uf.Uid = FilterKeep
	case "mine":
		switch mode {
		case FilterPurposePack:
			return uf, fmt.Errorf("filter UID cannot use 'mine' mode: only makes sense when unpacking")
		case FilterPurposeUnpack:
			uf.Uid = FilterMine
		}
	default:
		uf.Uid, err = strconv.Atoi(ff.Uid)
		if err != nil || uf.Uid < 0 {
			return uf, fmt.Errorf("filter UID must be one of 'keep', 'mine', or a positive int")
		}
	}

	// Parse GID.
	switch ff.Gid {
	case "":
		switch mode {
		case FilterPurposePack:
			uf.Gid = DefaultGid
		case FilterPurposeUnpack:
			uf.Gid = FilterMine
		}
	case "keep":
		uf.Gid = FilterKeep
	case "mine":
		switch mode {
		case FilterPurposePack:
			return uf, fmt.Errorf("filter GID cannot use 'mine' mode: only makes sense when unpacking")
		case FilterPurposeUnpack:
			uf.Gid = FilterMine
		}
	default:
		uf.Gid, err = strconv.Atoi(ff.Gid)
		if err != nil || uf.Gid < 0 {
			return uf, fmt.Errorf("filter GID must be one of 'keep', 'mine', or a positive int")
		}
	}

	// Parse time.
	switch {
	case ff.Mtime == "":
		switch mode {
		case FilterPurposePack:
			uf.Mtime = &DefaultMtime
		case FilterPurposeUnpack:
			uf.Mtime = nil // 'keep'
		}
	case ff.Mtime == "keep":
		uf.Mtime = nil
	case ff.Mtime[1] == '@':
		ut, err := strconv.Atoi(ff.Mtime[1:])
		if err != nil {
			return uf, fmt.Errorf("filter mtime parameter starting with '@' must be unix timestamp integer")
		}
		*uf.Mtime = time.Unix(int64(ut), 0)
	default:
		*uf.Mtime, err = time.Parse(time.RFC3339, ff.Mtime)
		if err != nil {
			return uf, fmt.Errorf("filter mtime parameter must be either 'keep', a unix timestamp integer beginning with '@', or an RFC3339 date string")
		}
	}

	// Parse sticky -- relatively simple.
	//  But despite being logically a bool, requires three-value logic to communicate "default".
	switch ff.Sticky {
	case "":
		switch mode {
		case FilterPurposePack:
			uf.Sticky = true
		case FilterPurposeUnpack:
			uf.Sticky = false
		}
	case "keep":
		uf.Sticky = true
	case "zero":
		uf.Sticky = false
	default:
		return uf, fmt.Errorf("sticky mode either 'keep' or 'zero'")
	}

	return
}

/*
	This is analogous to the serializable API FilesetFilters struct,
	but uses fields of more useful types instead of worrying about being serializable.

	Instances are produced by the `ProcessFilters()` function,
	which rejects any values that are out of range -- so
	code is free to presume fields only exist within their valid ranges when using this type.
*/
type FilesetFilters struct {
	Uid    int        // -1 for "keep", -2 for "mine"
	Gid    int        // -1 for "keep", -2 for "mine"
	Mtime  *time.Time // nil for "keep"
	Sticky bool
}
