package api

func (ref SubmoduleRef) Contextualize(parent SubmoduleRef) SubmoduleRef {
	if ref == "" {
		return parent
	}
	return SubmoduleRef(string(parent) + "." + string(ref))
}
func (ref SubmoduleStepRef) Contextualize(parent SubmoduleRef) SubmoduleStepRef {
	return SubmoduleStepRef{
		ref.SubmoduleRef.Contextualize(parent),
		ref.StepName,
	}
}
func (ref SubmoduleSlotRef) Contextualize(parent SubmoduleRef) SubmoduleSlotRef {
	return SubmoduleSlotRef{
		ref.SubmoduleRef.Contextualize(parent),
		ref.SlotRef,
	}
}
