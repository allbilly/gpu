package rdnaemu

import "gitlab.com/akita/navisim/rdnainsts"

// Decoder defines the interface that can convert instruction bytes to
// instructions
type Decoder interface {
	Decode(buf []byte) (*rdnainsts.Inst, error)
}
