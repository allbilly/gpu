package rdnaemu

import (
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/util/ca"
)

// InstEmuState is the interface used by the emulator to track the instruction
// execution status.
type InstEmuState interface {
	PID() ca.PID
	Inst() *rdnainsts.Inst
	Scratchpad() Scratchpad
}
