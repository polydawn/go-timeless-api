package api

import (
	"testing"

	. "github.com/warpfork/go-wish"
)

func TestNilRelationLexicalOrdering(t *testing.T) {
	basting := Module{Operations: map[StepName]StepUnion{
		"stepD": Operation{},
		"stepB": Operation{},
		"stepA": Operation{},
		"stepC": Operation{},
	}}
	order, err := OrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, []StepName{
		"stepA",
		"stepB",
		"stepC",
		"stepD",
	})
}

func TestFanoutLexicalOrdering(t *testing.T) {
	basting := Module{Operations: map[StepName]StepUnion{
		"stepD": Operation{Inputs: map[SlotReference]AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"stepB": Operation{Inputs: map[SlotReference]AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"stepA": Operation{Inputs: map[SlotReference]AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"stepC": Operation{Inputs: map[SlotReference]AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"step0": Operation{Outputs: map[SlotName]AbsPath{
			"theslot": "/",
		}},
	}}
	order, err := OrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, []StepName{
		"step0",
		"stepA",
		"stepB",
		"stepC",
		"stepD",
	})
}
