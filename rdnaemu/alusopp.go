package rdnaemu

import "log"

//nolint:gocyclo
func (u *ALUImpl) runSOPP(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0: // S_NOP
	// Do nothing
	case 2: // S_CBRANCH
		u.runSCBRANCH(state)
	case 4: // S_CBRANCH_SCC0
		u.runSCBRANCHSCC0(state)
	case 5: // S_CBRANCH_SCC1
		u.runSCBRANCHSCC1(state)
	case 6: // S_CBRANCH_VCCZ
		u.runSCBRANCHVCCZ(state)
	case 7: // S_CBRANCH_VCCNZ
		u.runSCBRANCHVCCNZ(state)
	case 8: // S_CBRANCH_EXECZ
		u.runSCBRANCHEXECZ(state)
	case 9: // S_CBRANCH_EXECNZ
		u.runSCBRANCHEXECNZ(state)
	case 12: // S_WAITCNT
		// Do nothing
	case 33: // S_CLAUSE
		// Do nothing
	case 35: //S_WAINTCNT_DEPCTR
		//Do noting
	default:
		log.Panicf("Opcode %d for SOPP format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSCBRANCH(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
}

func (u *ALUImpl) runSCBRANCHSCC0(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	if sp.SCC == 0 {
		sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
	}
}

func (u *ALUImpl) runSCBRANCHSCC1(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	if sp.SCC == 1 {
		sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
	}
}

func (u *ALUImpl) runSCBRANCHVCCZ(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	if sp.VCC == 0 {
		sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
	}
}

func (u *ALUImpl) runSCBRANCHVCCNZ(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	if sp.VCC != 0 {
		sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
	}
}

func (u *ALUImpl) runSCBRANCHEXECZ(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	if sp.EXEC == 0 {
		sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
	}
}

func (u *ALUImpl) runSCBRANCHEXECNZ(state InstEmuState) {
	sp := state.Scratchpad().AsSOPP()
	imm := asInt16(uint16(sp.IMM & 0xffff))
	if sp.EXEC != 0 {
		sp.PC = uint64(int64(sp.PC) + int64(imm)*4)
	}
}

func laneMasked(Exec uint64, laneID uint) bool {
	return Exec&(1<<laneID) > 0
}
