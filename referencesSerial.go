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

// Despite the fact this is documented as a union (and it is), we're not
// using a style of serialization here that Refmt has no explicit support for
// (we're screwing with strings in an alarmingly intricate way), so,
// we're doing it on our own in some transform funcs.
var ImportRef_AtlasEntry = atlas.BuildEntry((*ImportRef)(nil)).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(func(x ImportRef) (string, error) { return x.String(), nil })).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseImportRef)).
	Complete()

func ParseImportRef(x string) (ImportRef, error) {
	hunks := strings.SplitN(x, ":", 2)
	if len(hunks) == 1 {
		return nil, fmt.Errorf("valid import refs begin with 'catalog', 'parent', or 'ingest', followed by a colon and additional information.")
	}
	switch hunks[0] {
	case "catalog":
		itemRef, err := ParseItemRef(hunks[1])
		return ImportRef_Catalog(itemRef), err
	case "parent":
		slotRef, err := ParseSlotReference(hunks[1])
		return ImportRef_Parent(slotRef), err
	case "ingest":
		panic("TODO") // TODO
	default:
		return nil, fmt.Errorf("valid import refs begin with 'catalog', 'parent', or 'ingest', followed by a colon and additional information.")
	}
	return nil, nil
}
