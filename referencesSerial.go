package api

import (
	"fmt"
	"strings"

	"github.com/polydawn/refmt/obj/atlas"
)

var SlotRef_AtlasEntry = atlas.BuildEntry(SlotRef{}).Transform().
	TransformMarshal(atlas.MakeMarshalTransformFunc(func(x SlotRef) (string, error) { return x.String(), nil })).
	TransformUnmarshal(atlas.MakeUnmarshalTransformFunc(ParseSlotRef)).
	Complete()

func ParseSlotRef(x string) (SlotRef, error) {
	if x == "" {
		return SlotRef{}, fmt.Errorf("empty string is not a SlotRef")
	}
	ss := strings.SplitN(x, ".", 2)
	switch len(ss) {
	case 1:
		return SlotRef{"", SlotName(ss[0])}, nil
	case 2:
		return SlotRef{StepName(ss[0]), SlotName(ss[1])}, nil
	default:
		return SlotRef{}, fmt.Errorf("slot references can be of form 'x' or 'x.y'.")
	}
}

func (x SlotRef) String() string {
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
		slotRef, err := ParseSlotRef(hunks[1])
		return ImportRef_Parent(slotRef), err
	case "ingest":
		ss := strings.SplitN(hunks[1], ":", 2)
		if len(ss) != 2 {
			return nil, fmt.Errorf("valid ingest refs begin with 'ingest:thetool', where \"thetool\" may be a name of an ingest system (such as \"git\").")
		}
		return ImportRef_Ingest{ss[0], ss[1]}, nil
	default:
		return nil, fmt.Errorf("valid import refs begin with 'catalog', 'parent', or 'ingest', followed by a colon and additional information.")
	}
	return nil, nil
}

func (x SubmoduleStepRef) String() string {
	if x.SubmoduleRef == "" && x.StepName == "" {
		return ""
	} else if x.SubmoduleRef == "" {
		return string(x.StepName)
	} else {
		return string(x.SubmoduleRef) + "." + string(x.StepName)
	}
}

func (x SubmoduleSlotRef) String() string {
	if x.SubmoduleRef == "" && x.StepName == "" && x.SlotName == "" {
		return ""
	} else if x.SubmoduleRef == "" && x.StepName == "" {
		return string(x.SlotName)
	} else if x.StepName == "" {
		return string(x.SubmoduleRef) + "." + string(x.SlotName)
	} else if x.SubmoduleRef == "" {
		return string(x.StepName) + "." + string(x.SlotName)
	} else {
		return string(x.SubmoduleRef) + "." + string(x.StepName) + "." + string(x.SlotName)
	}
}
