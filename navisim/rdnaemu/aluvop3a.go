package rdnaemu

import (
	"log"
	"math"
	"strings"

	"gitlab.com/akita/navisim/bitops"
)

//nolint:gocyclo,funlen
func (u *ALUImpl) runVOP3A(state InstEmuState) {
	inst := state.Inst()
	//log.Printf("%s\n", inst.InstName)
	u.vop3aPreprocess(state)

	switch inst.Opcode {
	case 4:
		u.runVCMPGTF32VOP3A(state)
	case 132:
		u.runVCMPGTI32VOP3A(state)
	case 134:
		u.runVCMPGEI32VOP3A(state)
	case 194:
		u.runVCMPEQU32VOP3A(state)
	case 257: //256+1
		u.runVCNDMASKB32VOP3A(state)
	case 267: //256+11
		u.runVMULU32U24VOP3A(state)
	case 323:
		u.runVMADU32U24(state)
	case 361:
		u.runVMULOU32(state)
	case 362:
		u.runVMULHIU32(state)
	case 767:
		u.runVLSHLREVB64(state)
	case 769:
		u.runVASHRREVI64(state)
	case 838:
		u.runVLSHLADDU32(state)
	case 877:
		u.runVADD3U32(state)
	default:
		log.Panicf("Opcode %d for VOP3a format is not implemented", inst.Opcode)
	}
	u.vop3aPostprocess(state)
}

func (u *ALUImpl) vop3aPreprocess(state InstEmuState) {
	inst := state.Inst()

	if inst.Abs != 0 {
		u.vop3aPreProcessAbs(state)
	}

	if inst.Neg != 0 {
		u.vop3aPreProcessNeg(state)
	}
}

func (u *ALUImpl) vop3aPreProcessAbs(state InstEmuState) {
	inst := state.Inst()
	sp := state.Scratchpad().AsVOP3A()

	if strings.Contains(inst.InstName, "F32") ||
		strings.Contains(inst.InstName, "f32") {
		if inst.Abs&0x1 != 0 {
			for i := 0; i < 64; i++ {
				src0 := math.Float32frombits(uint32(sp.SRC0[i]))
				src0 = float32(math.Abs(float64(src0)))
				sp.SRC0[i] = uint64(math.Float32bits(src0))
			}
		}

		if inst.Abs&0x2 != 0 {
			for i := 0; i < 64; i++ {
				src1 := math.Float32frombits(uint32(sp.SRC1[i]))
				src1 = float32(math.Abs(float64(src1)))
				sp.SRC1[i] = uint64(math.Float32bits(src1))
			}
		}

		if inst.Abs&0x4 != 0 {
			for i := 0; i < 64; i++ {
				src2 := math.Float32frombits(uint32(sp.SRC2[i]))
				src2 = float32(math.Abs(float64(src2)))
				sp.SRC2[i] = uint64(math.Float32bits(src2))
			}
		}
	} else {
		log.Printf("Absolute operation for %s is not implemented.", inst.InstName)
	}
}

func (u *ALUImpl) vop3aPreProcessNeg(state InstEmuState) {
	inst := state.Inst()

	if strings.Contains(inst.InstName, "F64") ||
		strings.Contains(inst.InstName, "f64") {
		u.vop3aPreProcessF64Neg(state)
	} else if strings.Contains(inst.InstName, "F32") ||
		strings.Contains(inst.InstName, "f32") {
		u.vop3aPreProcessF32Neg(state)
	} else if strings.Contains(inst.InstName, "B32") ||
		strings.Contains(inst.InstName, "b32") {
		u.vop3aPreProcessB32Neg(state)
	} else {
		log.Printf("Negative operation for %s is not implemented.", inst.InstName)
	}
}

func (u *ALUImpl) vop3aPreProcessF64Neg(state InstEmuState) {
	inst := state.Inst()
	sp := state.Scratchpad().AsVOP3A()

	if inst.Neg&0x1 != 0 {
		for i := 0; i < 64; i++ {
			src0 := math.Float64frombits(sp.SRC0[i])
			src0 = src0 * (-1.0)
			sp.SRC0[i] = math.Float64bits(src0)
		}
	}

	if inst.Neg&0x2 != 0 {
		for i := 0; i < 64; i++ {
			src1 := math.Float64frombits(sp.SRC1[i])
			src1 = src1 * (-1.0)
			sp.SRC1[i] = math.Float64bits(src1)
		}
	}

	if inst.Neg&0x4 != 0 {
		for i := 0; i < 64; i++ {
			src2 := math.Float64frombits(sp.SRC2[i])
			src2 = src2 * (-1.0)
			sp.SRC2[i] = math.Float64bits(src2)
		}
	}
}

func (u *ALUImpl) vop3aPreProcessF32Neg(state InstEmuState) {
	inst := state.Inst()
	sp := state.Scratchpad().AsVOP3A()
	if inst.Neg&0x1 != 0 {
		for i := 0; i < 64; i++ {
			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src0 = src0 * (-1.0)
			sp.SRC0[i] = uint64(math.Float32bits(src0))
		}
	}

	if inst.Neg&0x2 != 0 {
		for i := 0; i < 64; i++ {
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			src1 = src1 * (-1.0)
			sp.SRC1[i] = uint64(math.Float32bits(src1))
		}
	}

	if inst.Neg&0x4 != 0 {
		for i := 0; i < 64; i++ {
			src2 := math.Float32frombits(uint32(sp.SRC2[i]))
			src2 = src2 * (-1.0)
			sp.SRC2[i] = uint64(math.Float32bits(src2))
		}
	}
}

func (u *ALUImpl) vop3aPreProcessB32Neg(state InstEmuState) {
	inst := state.Inst()
	sp := state.Scratchpad().AsVOP3A()
	if inst.Neg&0x1 != 0 {
		for i := 0; i < 64; i++ {
			src0 := asInt32(uint32(sp.SRC0[i]))
			src0 = src0 * (-1.0)
			sp.SRC0[i] = uint64(int32ToBits(src0))
		}
	}

	if inst.Neg&0x2 != 0 {
		for i := 0; i < 64; i++ {
			src1 := asInt32(uint32(sp.SRC1[i]))
			src1 = src1 * (-1.0)
			sp.SRC1[i] = uint64(int32ToBits(src1))
		}
	}

	if inst.Neg&0x4 != 0 {
		for i := 0; i < 64; i++ {
			src2 := asInt32(uint32(sp.SRC2[i]))
			src2 = src2 * (-1.0)
			sp.SRC2[i] = uint64(int32ToBits(src2))
		}
	}
}

func (u *ALUImpl) vop3aPostprocess(state InstEmuState) {
	inst := state.Inst()

	if inst.Omod != 0 {
		log.Panic("Output modifiers are not supported.")
	}
}

func (u *ALUImpl) runVCMPGTF32VOP3A(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()
	var i uint
	var src0, src1 float32
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		src0 = math.Float32frombits(uint32(sp.SRC0[i]))
		src1 = math.Float32frombits(uint32(sp.SRC1[i]))
		if src0 > src1 {
			sp.DST[0] |= 1 << i
		}
	}
}

func (u *ALUImpl) runVCMPGTI32VOP3A(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		src0 := asInt32(uint32(sp.SRC0[i]))
		src1 := asInt32(uint32(sp.SRC1[i]))
		if src0 > src1 {
			sp.DST[0] |= 1 << i
		}
	}
}

func (u *ALUImpl) runVCMPGEI32VOP3A(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		src0 := asInt32(uint32(sp.SRC0[i]))
		src1 := asInt32(uint32(sp.SRC1[i]))
		if src0 >= src1 {
			sp.DST[0] |= 1 << i
		}
	}
}
func (u *ALUImpl) runVCMPEQU32VOP3A(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		src0 := uint32(sp.SRC0[i])
		src1 := uint32(sp.SRC1[i])
		if src0 == src1 {
			sp.DST[0] |= 1 << i
		}
	}
}

func (u *ALUImpl) runVASHRREVI64(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		sp.DST[i] = int64ToBits(asInt64(sp.SRC1[i]) >> sp.SRC0[i])
	}
}
func (u *ALUImpl) runVMULU32U24VOP3A(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		src0 := uint32(bitops.ExtractBitsFromU64(sp.SRC0[i], 0, 23))
		src1 := uint32(bitops.ExtractBitsFromU64(sp.SRC1[i], 0, 23))
		sp.DST[i] = uint64(src0 * src1)
	}
}

func (u *ALUImpl) runVCNDMASKB32VOP3A(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		if (sp.SRC2[i] & (1 << i)) > 0 {
			sp.DST[i] = sp.SRC1[i]
		} else {
			sp.DST[i] = sp.SRC0[i]
		}
	}
}

func (u *ALUImpl) runVLSHLREVB64(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		sp.DST[i] = sp.SRC1[i] << sp.SRC0[i]
	}
}

func (u *ALUImpl) runVMADU32U24(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		src0 := uint32(bitops.ExtractBitsFromU64(sp.SRC0[i], 0, 23))
		src1 := uint32(bitops.ExtractBitsFromU64(sp.SRC1[i], 0, 23))
		src2 := uint32(sp.SRC2[i])
		sp.DST[i] = uint64(src0*src1 + src2)
	}
}

func (u *ALUImpl) runVMULOU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		sp.DST[i] = (sp.SRC0[i] * sp.SRC1[i])
	}
}

func (u *ALUImpl) runVMULHIU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		sp.DST[i] = (sp.SRC0[i] * sp.SRC1[i]) >> 32
	}
}

func (u *ALUImpl) runVLSHLADDU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		sp.DST[i] = (sp.SRC0[i] << sp.SRC1[i]) + sp.SRC2[i]
	}
}

func (u *ALUImpl) runVADD3U32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP3A()

	var i uint
	for i = 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		sp.DST[i] = sp.SRC0[i] + sp.SRC1[i] + sp.SRC2[i]
	}
}
