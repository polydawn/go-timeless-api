package api

import (
	"fmt"
	"strings"

	"github.com/polydawn/refmt/obj/atlas"
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

var Atlas_Catalog = atlas.MustBuild(
	Lineage_AtlasEntry,
	Release_AtlasEntry,
	WareID_AtlasEntry,
)

var ItemRef_AtlasEntry = atlas.BuildEntry(ItemRef{}).StructMap().Autogenerate().Complete()
var Lineage_AtlasEntry = atlas.BuildEntry(Lineage{}).StructMap().Autogenerate().Complete()
var Release_AtlasEntry = atlas.BuildEntry(Release{}).StructMap().Autogenerate().Complete()
