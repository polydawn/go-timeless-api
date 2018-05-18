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

func TestFanInLexicalOrdering(t *testing.T) {
	basting := Module{Operations: map[StepName]StepUnion{
		"stepD": Operation{Outputs: map[SlotName]AbsPath{"theslot": "/"}},
		"stepB": Operation{Outputs: map[SlotName]AbsPath{"theslot": "/"}},
		"stepA": Operation{Outputs: map[SlotName]AbsPath{"theslot": "/"}},
		"stepC": Operation{Outputs: map[SlotName]AbsPath{"theslot": "/"}},
		"step9": Operation{Inputs: map[SlotReference]AbsPath{
			{"stepA", "theslot"}: "/a",
			{"stepB", "theslot"}: "/b",
			{"stepC", "theslot"}: "/c",
			{"stepD", "theslot"}: "/d",
		}},
	}}
	order, err := OrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, []StepName{
		"stepA",
		"stepB",
		"stepC",
		"stepD",
		"step9",
	})
}

func TestSimpleLinearOrdering(t *testing.T) {
	basting := Module{Operations: map[StepName]StepUnion{
		"stepA": Operation{
			Outputs: map[SlotName]AbsPath{"aslot": "/"},
		},
		"stepB": Operation{
			Inputs:  map[SlotReference]AbsPath{{"stepA", "aslot"}: "/"},
			Outputs: map[SlotName]AbsPath{"xslot": "/"},
		},
		"stepD": Operation{
			Inputs:  map[SlotReference]AbsPath{{"stepB", "xslot"}: "/"},
			Outputs: map[SlotName]AbsPath{"xslot": "/"},
		},
		"stepC": Operation{
			Inputs: map[SlotReference]AbsPath{{"stepD", "xslot"}: "/"},
		},
	}}
	order, err := OrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, []StepName{
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
	basting := Module{Operations: map[StepName]StepUnion{
		"stepA": Operation{
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepB": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepA", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepC": Operation{
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepD": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepC", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepE": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepB", "slot"}: "/",
				{"stepD", "slot"}: "/1",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepF": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepD", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepG": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepF", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepH": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepE", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepI": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepF", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepJ": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepF", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepK": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepE", "slot"}: "/",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepL": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepG", "slot"}: "/",
				{"stepK", "slot"}: "/1",
				{"stepH", "slot"}: "/2",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepM": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepI", "slot"}: "/",
				{"stepJ", "slot"}: "/1",
			},
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
	}}
	order, err := OrderSteps(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, []StepName{
		"stepA", "stepB", "stepC", "stepD", "stepE",
		"stepF", "stepG", "stepH", "stepI", "stepJ",
		"stepK", "stepL", "stepM",
	})
}

func TestDeepSubmoduleOrdering(t *testing.T) {
	basting := Module{Operations: map[StepName]StepUnion{
		"stepFoo": Operation{
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepBar": Operation{
			Outputs: map[SlotName]AbsPath{"slot": "/"},
		},
		"stepSub": Module{
			Imports: map[SlotName]ImportPattern{
				"subx": {"parent", "stepFoo.slot"},
			},
			Operations: map[StepName]StepUnion{
				"deeper": Module{
					Imports: map[SlotName]ImportPattern{
						"suby": {"parent", "subx"},
					},
					Operations: map[StepName]StepUnion{
						"rlydeep": Operation{
							Outputs: map[SlotName]AbsPath{"slot": "/"},
						},
					},
					Exports: map[ItemName]SlotReference{
						"zowslot": {"midstep", "slot"},
					},
				},
				"midstep": Operation{
					Inputs:  map[SlotReference]AbsPath{{"deeper", "zowslot"}: "/"},
					Outputs: map[SlotName]AbsPath{"slot": "/"},
				},
			},
			Exports: map[ItemName]SlotReference{
				"wowslot": {"midstep", "slot"},
			},
		},
		"stepWub": Operation{
			Inputs: map[SlotReference]AbsPath{
				{"stepSub", "wowslot"}: "/",
			},
		},
	}}
	order, err := OrderStepsDeep(basting)
	Wish(t, err, ShouldEqual, nil)
	Wish(t, order, ShouldEqual, []SubmoduleStepReference{
		SubmoduleStepReference{"", "stepBar"},
		SubmoduleStepReference{"", "stepFoo"},
		SubmoduleStepReference{"stepSub.deeper", "rlydeep"},
		SubmoduleStepReference{"stepSub", "midstep"},
		SubmoduleStepReference{"", "stepWub"},
	})
}
