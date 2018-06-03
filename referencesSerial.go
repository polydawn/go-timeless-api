package api

import (
	"fmt"
	"strings"

	"github.com/polydawn/refmt/obj/atlas"
)

var SlotReference_AtlasEntry = atlas.BuildEntry(SlotReference{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(func(x SlotReference) (string, error) { return x.String(), nil })).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseSlotReference)).
	Complete()

func ParseSlotReference(x string) (SlotReference, error) {
	if x == "" {
		return SlotReference{}, fmt.Errorf("empty string is not a SlotReference")
	}
	ss := strings.SplitN(x, ".", 2)
	switch len(ss) {
	case 1:
		return SlotReference{"", SlotName(ss[0])}, nil
	case 2:
		return SlotReference{StepName(ss[0]), SlotName(ss[1])}, nil
	default:
		return SlotReference{}, fmt.Errorf("slot references can be of form 'x' or 'x.y'.")
	}
}

func (x SlotReference) String() string {
	if x.StepName == "" && x.SlotName == "" {
		return ""
	} else if x.StepName == "" {
		return string(x.SlotName)
	} else {
		return string(x.StepName) + "." + string(x.SlotName)
	}
}
