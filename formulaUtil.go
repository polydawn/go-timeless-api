package api

import (
	"github.com/polydawn/refmt"
)

func (f Formula) Clone() (f2 Formula) {
	refmt.MustCloneAtlased(f, &f2, Atlas_Formula)
	return
}
