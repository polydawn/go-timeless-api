package api

func (ref SubmoduleRef) Contextualize(parent StepName) SubmoduleRef {
	if ref == "" {
		return SubmoduleRef(parent)
	}
	return SubmoduleRef(string(parent) + "." + string(ref))
}
func (ref SubmoduleStepRef) Contextualize(parent StepName) SubmoduleStepRef {
	return SubmoduleStepRef{
		ref.SubmoduleRef.Contextualize(parent),
		ref.StepName,
	}
}
func (ref SubmoduleSlotRef) Contextualize(parent StepName) SubmoduleSlotRef {
	return SubmoduleSlotRef{
		ref.SubmoduleRef.Contextualize(parent),
		ref.SlotRef,
	}
}
