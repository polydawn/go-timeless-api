package api

import (
	"bytes"
	"testing"

	"github.com/polydawn/refmt"
	"github.com/polydawn/refmt/json"
	. "github.com/warpfork/go-wish"
)

func TestModuleSerialization(t *testing.T) {
	t.Run("zero module should roundtrip", func(t *testing.T) {
		obj := Module{}
		canon := `{"imports":null,"steps":null}`

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{}, obj, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := Module{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
		t.Run("unmarshal blank", func(t *testing.T) {
			targ := Module{}
			canon := `{}`
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
	t.Run("exhilerating module should roundtrip", func(t *testing.T) {
		obj := Module{Operations: map[StepName]StepUnion{
			"stepFoo": Operation{
				Outputs: map[SlotName]AbsPath{"slot": "/"},
			},
			"stepBar": Operation{
				Outputs: map[SlotName]AbsPath{"slot": "/"},
			},
			"stepSub": Module{
				Imports: map[SlotName]ImportRef{
					"subx": ImportRef_Parent{"stepFoo", "slot"},
				},
				Operations: map[StepName]StepUnion{
					"deeper": Module{
						Imports: map[SlotName]ImportRef{
							"suby": ImportRef_Parent{"", "subx"},
						},
						Operations: map[StepName]StepUnion{
							"rlydeep": Operation{
								Outputs: map[SlotName]AbsPath{"slot": "/"},
							},
						},
						Exports: map[ItemName]SlotReference{
							"zowslot": {"rlydeep", "slot"},
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
		canon := Dedent(`
			{
				"imports": null,
				"steps": {
					"stepBar": {
						"operation": {
							"inputs": null,
							"action": {},
							"outputs": {
								"slot": "/"
							}
						}
					},
					"stepFoo": {
						"operation": {
							"inputs": null,
							"action": {},
							"outputs": {
								"slot": "/"
							}
						}
					},
					"stepSub": {
						"module": {
							"imports": {
								"subx": "parent:stepFoo.slot"
							},
							"steps": {
								"deeper": {
									"module": {
										"imports": {
											"suby": "parent:subx"
										},
										"steps": {
											"rlydeep": {
												"operation": {
													"inputs": null,
													"action": {},
													"outputs": {
														"slot": "/"
													}
												}
											}
										},
										"exports": {
											"zowslot": "rlydeep.slot"
										}
									}
								},
								"midstep": {
									"operation": {
										"inputs": {
											"deeper.zowslot": "/"
										},
										"action": {},
										"outputs": {
											"slot": "/"
										}
									}
								}
							},
							"exports": {
								"wowslot": "midstep.slot"
							}
						}
					},
					"stepWub": {
						"operation": {
							"inputs": {
								"stepSub.wowslot": "/"
							},
							"action": {}
						}
					}
				}
			}
		`)

		t.Run("marshal", func(t *testing.T) {
			bs, err := refmt.MarshalAtlased(json.EncodeOptions{Line: []byte{'\n'}, Indent: []byte{'\t'}}, obj, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, string(bs), ShouldEqual, canon)
		})
		t.Run("unmarshal", func(t *testing.T) {
			targ := Module{}
			err := refmt.UnmarshalAtlased(json.DecodeOptions{}, bytes.NewBufferString(canon).Bytes(), &targ, Atlas_Module)
			Wish(t, err, ShouldEqual, nil)
			Wish(t, targ, ShouldEqual, obj)
		})
	})
}
