package api

func (ref SubmoduleReference) Contextualize(parent StepName) SubmoduleReference {
	if ref == "" {
		return SubmoduleReference(parent)
	}
	return SubmoduleReference(string(parent) + "." + string(ref))
}
func (ref SubmoduleStepReference) Contextualize(parent StepName) SubmoduleStepReference {
	return SubmoduleStepReference{
		ref.SubmoduleReference.Contextualize(parent),
		ref.StepName,
	}
}
