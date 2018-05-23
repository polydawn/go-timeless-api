package api

import (
	"fmt"
	"strings"
)

func ParseItemRef(x string) (v ItemRef, err error) {
	ss := strings.Split(x, ":")
	switch len(ss) {
	case 3:
		v.ItemName = ItemName(ss[2])
		fallthrough
	case 2:
		v.ReleaseName = ReleaseName(ss[1])
		fallthrough
	case 1:
		v.ModuleName = ModuleName(ss[0])
		return
	default:
		return ItemRef{}, fmt.Errorf("ReleaseItemIDs are a colon-separated three-tuple; no more than two colons may appear!")
	}
}

func (x ItemRef) String() string {
	switch {
	case x.ModuleName == "":
		return ""
	case x.ReleaseName == "":
		return string(x.ModuleName)
	case x.ItemName == "":
		return string(x.ModuleName) + ":" + string(x.ReleaseName)
	default:
		return string(x.ModuleName) + ":" + string(x.ReleaseName) + ":" + string(x.ItemName)
	}
}
