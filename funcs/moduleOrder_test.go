package funcs

import (
	"testing"

	. "github.com/warpfork/go-wish"

	"go.polydawn.net/go-timeless-api"
)

func TestNilRelationLexicalOrdering(t *testing.T) {
	basting := api.Module{Steps: map[api.StepName]api.StepUnion{
		"stepD": api.Operation{},
		"stepB": api.Operation{},
		"stepA": api.Operation{},
		"stepC": api.Operation{},
	}}
	order, err := ModuleOrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, StepList{
		"stepA",
		"stepB",
		"stepC",
		"stepD",
	})
}

func TestFanoutLexicalOrdering(t *testing.T) {
	basting := api.Module{Steps: map[api.StepName]api.StepUnion{
		"stepD": api.Operation{Inputs: map[api.SlotRef]api.AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"stepB": api.Operation{Inputs: map[api.SlotRef]api.AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"stepA": api.Operation{Inputs: map[api.SlotRef]api.AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"stepC": api.Operation{Inputs: map[api.SlotRef]api.AbsPath{
			{"step0", "theslot"}: "/",
		}},
		"step0": api.Operation{Outputs: map[api.SlotName]api.AbsPath{
			"theslot": "/",
		}},
	}}
	order, err := ModuleOrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, StepList{
		"step0",
		"stepA",
		"stepB",
		"stepC",
		"stepD",
	})
}

func TestFanInLexicalOrdering(t *testing.T) {
	basting := api.Module{Steps: map[api.StepName]api.StepUnion{
		"stepD": api.Operation{Outputs: map[api.SlotName]api.AbsPath{"theslot": "/"}},
		"stepB": api.Operation{Outputs: map[api.SlotName]api.AbsPath{"theslot": "/"}},
		"stepA": api.Operation{Outputs: map[api.SlotName]api.AbsPath{"theslot": "/"}},
		"stepC": api.Operation{Outputs: map[api.SlotName]api.AbsPath{"theslot": "/"}},
		"step9": api.Operation{Inputs: map[api.SlotRef]api.AbsPath{
			{"stepA", "theslot"}: "/a",
			{"stepB", "theslot"}: "/b",
			{"stepC", "theslot"}: "/c",
			{"stepD", "theslot"}: "/d",
		}},
	}}
	order, err := ModuleOrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, StepList{
		"stepA",
		"stepB",
		"stepC",
		"stepD",
		"step9",
	})
}

func TestSimpleLinearOrdering(t *testing.T) {
	basting := api.Module{Steps: map[api.StepName]api.StepUnion{
		"stepA": api.Operation{
			Outputs: map[api.SlotName]api.AbsPath{"aslot": "/"},
		},
		"stepB": api.Operation{
			Inputs:  map[api.SlotRef]api.AbsPath{{"stepA", "aslot"}: "/"},
			Outputs: map[api.SlotName]api.AbsPath{"xslot": "/"},
		},
		"stepD": api.Operation{
			Inputs:  map[api.SlotRef]api.AbsPath{{"stepB", "xslot"}: "/"},
			Outputs: map[api.SlotName]api.AbsPath{"xslot": "/"},
		},
		"stepC": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{{"stepD", "xslot"}: "/"},
		},
	}}
	order, err := ModuleOrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, StepList{
		"stepA",
		"stepB",
		"stepD",
		"stepC",
	})
}

func TestComplexOrdering(t *testing.T) {
	/*
	               /------------> K --\
	               |                   \
	  A --> B -----E ------> H --------> L
	              /                     /
	    C --> D -----F --> G ----------/
	                 |
	                 \------> I----------> M
	                 |                    /
	                 \--------> J -------/
	*/
	basting := api.Module{Steps: map[api.StepName]api.StepUnion{
		"stepA": api.Operation{
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepB": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepA", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepC": api.Operation{
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepD": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepC", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepE": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepB", "slot"}: "/",
				{"stepD", "slot"}: "/1",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepF": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepD", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepG": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepF", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepH": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepE", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepI": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepF", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepJ": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepF", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepK": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepE", "slot"}: "/",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepL": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepG", "slot"}: "/",
				{"stepK", "slot"}: "/1",
				{"stepH", "slot"}: "/2",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepM": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepI", "slot"}: "/",
				{"stepJ", "slot"}: "/1",
			},
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
	}}
	order, err := ModuleOrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, StepList{
		"stepA", "stepB", "stepC", "stepD", "stepE",
		"stepF", "stepG", "stepH", "stepI", "stepJ",
		"stepK", "stepL", "stepM",
	})
}

func TestDeepSubmoduleOrdering(t *testing.T) {
	basting := api.Module{Steps: map[api.StepName]api.StepUnion{
		"stepFoo": api.Operation{
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepBar": api.Operation{
			Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
		},
		"stepSub": api.Module{
			Imports: map[api.SlotName]api.ImportRef{
				"subx": api.ImportRef_Parent{"stepFoo", "slot"},
			},
			Steps: map[api.StepName]api.StepUnion{
				"deeper": api.Module{
					Imports: map[api.SlotName]api.ImportRef{
						"suby": api.ImportRef_Parent{"", "subx"},
					},
					Steps: map[api.StepName]api.StepUnion{
						"rlydeep": api.Operation{
							Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
						},
					},
					Exports: map[api.ItemName]api.SlotRef{
						"zowslot": {"rlydeep", "slot"},
					},
				},
				"midstep": api.Operation{
					Inputs:  map[api.SlotRef]api.AbsPath{{"deeper", "zowslot"}: "/"},
					Outputs: map[api.SlotName]api.AbsPath{"slot": "/"},
				},
			},
			Exports: map[api.ItemName]api.SlotRef{
				"wowslot": {"midstep", "slot"},
			},
		},
		"stepWub": api.Operation{
			Inputs: map[api.SlotRef]api.AbsPath{
				{"stepSub", "wowslot"}: "/",
			},
		},
	}}
	order, err := ModuleOrderStepsDeep(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, StepTree{
		{"", "stepBar"},
		{"", "stepFoo"},
		{"", "stepSub"},
		{"stepSub", "deeper"},
		{"stepSub.deeper", "rlydeep"},
		{"stepSub", "midstep"},
		{"", "stepWub"},
	})
}

func TestStepTreeDetach(t *testing.T) {
	t1 := StepTree{
		{"", "stepBar"},
		{"", "stepFoo"},
		{"", "stepSub"},
		{"stepSub", "deeper"},
		{"stepSub.deeper", "rlydeep"},
		{"stepSub", "midstep"},
		{"", "stepWub"},
	}
	t2 := t1.DetachSubtree("stepSub")
	Wish(t, t2, ShouldEqual, StepTree{
		{"", "deeper"},
		{"deeper", "rlydeep"},
		{"", "midstep"},
	})
	t3 := t2.DetachSubtree("deeper")
	Wish(t, t3, ShouldEqual, StepTree{
		{"", "rlydeep"},
	})
}

// TODO referential integrity checks: exports must actually refer to local outputs or imports
