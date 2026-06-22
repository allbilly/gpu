package rdnaemu

import (
	"log"
)

func (u *ALUImpl) runVOP3B(state InstEmuState) {
	inst := state.Inst()

	u.vop3aPreprocess(state)

	switch inst.Opcode {
	case 296:
		u.runVADDCOCIU32VOP3B(state)
	case 783:
		u.runVADDCOU32(state)
	default:
		log.Panicf("Opcode %d for VOP3b format is not implemented", inst.Opcode)
	}

	u.vop3aPostprocess(state)
}

func (u *ALUImpl) runVADDCOU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3B()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		sp.DST[i] = sp.SRC0[i] + sp.SRC1[i] + sp.SRC2[i]
		if sp.DST[i] > 0xffffffff {
			sp.SDST |= 1 << i
			sp.DST[i] &= 0xffffffff
		}
	}
}
func (u *ALUImpl) runVADDCOCIU32VOP3B(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3B()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		sp.DST[i] = sp.SRC0[i] + sp.SRC1[i] + ((sp.SRC2[i] & (1 << i)) >> i)
		carry := uint64(0)
		if sp.DST[i] > 0xffffffff {
			carry = 1
		}
		sp.SDST |= carry << i
		sp.DST[i] &= 0xffffffff
	}
}
