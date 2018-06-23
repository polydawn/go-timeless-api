package api

import (
	"strings"
)

// First returns the first StepName component of the SubmoduleRef.
// The empty string is returned if this SubmoduleRef is itself zero.
func (ref SubmoduleRef) First() StepName {
	i := strings.IndexRune(string(ref), '.')
	if i > 0 {
		return StepName(ref[:i])
	}
	return StepName(ref)
}

// Child appends the stepname to this ref.
// Think of it as leaving breadcrumbs behind as you zoom in
// ('Child' and 'Decontextualize' often come in pairs.)
func (ref SubmoduleRef) Child(child StepName) SubmoduleRef {
	if ref == "" {
		return SubmoduleRef(child)
	}
	return SubmoduleRef(string(ref) + "." + string(child))
}

// Contextualize prepends a set of step references to this ref.
// Think of it as zooming out.
func (ref SubmoduleRef) Contextualize(parent SubmoduleRef) SubmoduleRef {
	if ref == "" {
		return parent
	}
	return SubmoduleRef(string(parent) + "." + string(ref))
}

// Decontextualize strips the first stepName from the front of the ref.
// Think of it as zooming in.
func (ref SubmoduleRef) Decontextualize() SubmoduleRef {
	i := strings.IndexRune(string(ref), '.')
	if i > 0 {
		return ref[i+1:]
	}
	return ""
}

// Contextualize prepends a set of step references to this ref.
func (ref SubmoduleStepRef) Contextualize(parent SubmoduleRef) SubmoduleStepRef {
	return SubmoduleStepRef{
		ref.SubmoduleRef.Contextualize(parent),
		ref.StepName,
	}
}

// Decontextualize strips the first stepName from the front of the ref.
// Think of it as zooming in.
func (ref SubmoduleStepRef) Decontextualize() SubmoduleStepRef {
	return SubmoduleStepRef{ref.SubmoduleRef.Decontextualize(), ref.StepName}
}

// Contextualize prepends a set of step references to this ref.
func (ref SubmoduleSlotRef) Contextualize(parent SubmoduleRef) SubmoduleSlotRef {
	return SubmoduleSlotRef{
		ref.SubmoduleRef.Contextualize(parent),
		ref.SlotRef,
	}
}

// Decontextualize strips the first stepName from the front of the ref.
// Think of it as zooming in.
func (ref SubmoduleSlotRef) Decontextualize() SubmoduleSlotRef {
	return SubmoduleSlotRef{
		ref.SubmoduleRef.Decontextualize(),
		ref.SlotRef,
	}
}
