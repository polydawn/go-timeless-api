package api

import (
	"fmt"
	"strings"

	"github.com/polydawn/refmt/obj/atlas"
)

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
