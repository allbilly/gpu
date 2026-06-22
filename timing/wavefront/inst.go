// Package wavefront defines concepts related to a wavefront.
package wavefront

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/navisim/rdnainsts"
)

// Inst in the timing package is a wrapper of the rdnainsts.Inst.
type Inst struct {
	*rdnainsts.Inst

	ID string
}

// NewInst creates a newly created Inst
func NewInst(raw *rdnainsts.Inst) *Inst {
	i := new(Inst)
	i.Inst = raw

	i.ID = akita.GetIDGenerator().Generate()

	return i
}
