package rdnaemu

import (
	"log"
	"math"
)

//nolint:gocyclo,funlen
func (u *ALUImpl) runSOP2(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSADDU32(state)
	case 2:
		u.runSADDI32(state)
	case 3:
		u.runSSUBI32(state)
	case 4:
		u.runSADDCU32(state)
	case 7:
		u.runSMINU32(state)
	case 14:
		u.runSANDB32(state)
	case 15:
		u.runSANDB64(state)
	case 16:
		u.runSORB32(state)
	case 18:
		u.runSXORB32(state)
	case 20:
		u.runSANDN2B32(state)
	case 30:
		u.runSLSHLB32(state)
	case 31:
		u.runSLSHLB64(state)
	case 32:
		u.runSLSHRB32(state)
	case 34:
		u.runSASHRI32(state)
	case 38:
		u.runSMULI32(state)
	default:
		log.Panicf("Opcode %d for SOP2 format is not implemented", inst.Opcode)
	}
}

func (u *ALUImpl) runSADDU32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	src0 := uint32(sp.SRC0)
	src1 := uint32(sp.SRC1)
	dst := src0 + src1

	if src0 > math.MaxUint32-src1 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
	sp.DST = uint64(dst)
}

func (u *ALUImpl) runSSUBI32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	src0 := asInt32(uint32(sp.SRC0))
	src1 := asInt32(uint32(sp.SRC1))
	dst := src0 - src1

	if src1 > 0 && dst > src0 {
		sp.SCC = 1
	} else if src1 < 0 && dst < src0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}

	sp.DST = uint64(int32ToBits(dst))
}

func (u *ALUImpl) runSADDCU32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	src0 := uint32(sp.SRC0)
	src1 := uint32(sp.SRC1)
	scc := uint32(sp.SCC)
	dst := src0 + src1 + scc

	if src0 > math.MaxUint32-src1-scc {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
	sp.DST = uint64(dst)
}

func (u *ALUImpl) runSMINU32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	if sp.SRC0 < sp.SRC1 {
		sp.DST = sp.SRC0
		sp.SCC = 1
	} else {
		sp.DST = sp.SRC1
	}
}

func (u *ALUImpl) runSADDI32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	src0 := asInt32(uint32(sp.SRC0))
	src1 := asInt32(uint32(sp.SRC1))
	dst := src0 + src1

	if src0 > math.MaxInt32-src1 {
		sp.SCC = 1
	}
	if src0 < math.MinInt32-src1 {
		sp.SCC = 1
	}

	sp.DST = uint64(int32ToBits(dst))
}

func (u *ALUImpl) runSANDB32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	sp.DST = sp.SRC0 & sp.SRC1
	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}
func (u *ALUImpl) runSANDB64(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	sp.DST = sp.SRC0 & sp.SRC1
	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}

func (u *ALUImpl) runSLSHLB32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	src0 := uint32(sp.SRC0)
	src1 := uint8(sp.SRC1)
	dst := src0 << (src1 & 0x1f)

	sp.DST = uint64(dst)
}

func (u *ALUImpl) runSLSHRB32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	sp.DST = sp.SRC0 >> (sp.SRC1 & 0x1f)

	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}

func (u *ALUImpl) runSLSHLB64(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	src0 := sp.SRC0
	src1 := uint8(sp.SRC1)
	dst := src0 << (src1 & 0x3f)

	sp.DST = dst
}

func (u *ALUImpl) runSORB32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	sp.DST = sp.SRC0 | sp.SRC1
	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}

func (u *ALUImpl) runSMULI32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	src0 := asInt32(uint32(sp.SRC0))
	src1 := asInt32(uint32(sp.SRC1))
	dst := src0 * src1

	sp.DST = uint64(int32ToBits(dst))

	if src0 != 0 && dst/src0 != src1 {
		sp.SCC = 1
	}
}
func (u *ALUImpl) runSASHRI32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	src0 := asInt32(uint32(sp.SRC0))
	src1 := uint8(sp.SRC1)
	dst := src0 >> src1

	sp.DST = uint64(int32ToBits(dst))

	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}

func (u *ALUImpl) runSXORB32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()
	src0 := uint32(sp.SRC0)
	src1 := uint32(sp.SRC1)
	dst := src0 ^ src1

	sp.DST = uint64(dst)

	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}
func (u *ALUImpl) runSANDN2B32(state InstEmuState) {
	sp := state.Scratchpad().AsSOP2()

	sp.DST = sp.SRC0 &^ sp.SRC1
	if sp.DST != 0 {
		sp.SCC = 1
	} else {
		sp.SCC = 0
	}
}
