package rdnaemu

import (
	"log"
)

func (u *ALUImpl) runSOPK(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 23: // S_WAITCNT_VSCNT
		// Do nothing
	default:
		log.Panicf("Opcode %d for SOPK format is not implemented", inst.Opcode)
	}
}
