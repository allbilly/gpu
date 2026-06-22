package rdnaemu

import (
	"log"
	"math"
)

//nolint:gocyclo,funlen
func (u *ALUImpl) runVOP2(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 1:
		u.runVCNDMASKB32(state)
	case 3:
		u.runVADDF32(state)
	case 4:
		u.runVSUBF32(state)
	case 5:
		u.runVSUBREVF32(state)
	case 8:
		u.runVMULF32(state)
	case 11:
		u.runVMULU32U24(state)
	case 15:
		u.runVMINF32(state)
	case 16:
		u.runVMAXF32(state)
	case 17:
		u.runVMINI32(state)
	case 18:
		u.runVMAXI32(state)
	case 28:
		u.runVORB32(state)
	case 24:
		u.runVASHRREVI32(state)
	case 26:
		u.runVLSHLREVB32(state)
	case 29:
		u.runVXORB32(state)
	case 37:
		u.runVADDNCU32(state)
	case 38:
		u.runVSUBNCU32(state)
	case 39:
		u.runVSUBREVNCU32(state)
	case 40:
		u.runVADDCOCIU32(state)
	case 43:
		u.runVFMACF32(state)

	default:
		log.Panicf("Opcode %d for VOP2 format (%s) is not implemented",
			inst.Opcode, inst.String(nil))
	}
}

func (u *ALUImpl) runVADDF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			dst := src0 + src1
			sp.DST[i] = uint64(math.Float32bits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode  %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVSUBF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			dst := src0 - src1
			sp.DST[i] = uint64(math.Float32bits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode  %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVSUBREVF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			dst := src1 - src0
			sp.DST[i] = uint64(math.Float32bits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode  %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVMULU32U24(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()

	for i := 0; i < 64; i++ {
		if !laneMasked(sp.EXEC, uint(i)) {
			continue
		}

		src0 := (uint32(sp.SRC0[i]) << 8) >> 8
		src1 := (uint32(sp.SRC1[i]) << 8) >> 8
		dst := src0 * src1
		sp.DST[i] = uint64(dst)
	}
}

func (u *ALUImpl) runVCNDMASKB32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			if (sp.VCC & (1 << i)) > 0 {
				sp.DST[i] = sp.SRC1[i]
			} else {
				sp.DST[i] = sp.SRC0[i]
			}
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVMULF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			dst := src0 * src1
			sp.DST[i] = uint64(math.Float32bits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVLSHLREVB32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := uint32(sp.SRC1[i])
			dst := src1 << (src0 & 0x1f)
			sp.DST[i] = uint64(dst)
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode  %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVASHRREVI32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := int32(sp.SRC1[i])
			dst := src1 >> (src0 & 0x1f)
			sp.DST[i] = uint64(dst)
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode  %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVMINF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			dst := src0
			if src1 < src0 {
				dst = src1
			}

			sp.DST[i] = uint64(math.Float32bits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVMAXF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := math.Float32frombits(uint32(sp.SRC0[i]))
			src1 := math.Float32frombits(uint32(sp.SRC1[i]))
			dst := src0
			if src1 > src0 {
				dst = src1
			}

			sp.DST[i] = uint64(math.Float32bits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVADDNCU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := uint32(sp.SRC1[i])
			dst := src0 + src1
			sp.DST[i] = uint64(dst)
		}
	} else {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := u.sdwaSrcSelect(uint32(sp.SRC0[i]), inst.Src0Sel)
			src1 := u.sdwaSrcSelect(uint32(sp.SRC1[i]), inst.Src1Sel)
			dst := src0 + src1
			dst = u.sdwaDstSelect(uint32(sp.DST[i]), dst,
				inst.DstSel, inst.DstUnused)
			sp.DST[i] = uint64(dst)
		}
	}
}
func (u *ALUImpl) runVSUBNCU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := uint32(sp.SRC1[i])
			dst := src0 - src1
			sp.DST[i] = uint64(dst)
		}
	} else {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := u.sdwaSrcSelect(uint32(sp.SRC0[i]), inst.Src0Sel)
			src1 := u.sdwaSrcSelect(uint32(sp.SRC1[i]), inst.Src1Sel)
			dst := src0 - src1
			dst = u.sdwaDstSelect(uint32(sp.DST[i]), dst,
				inst.DstSel, inst.DstUnused)
			sp.DST[i] = uint64(dst)
		}
	}
}

func (u *ALUImpl) runVSUBREVNCU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := uint32(sp.SRC1[i])
			dst := src1 - src0
			sp.DST[i] = uint64(dst)
		}
	} else {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := u.sdwaSrcSelect(uint32(sp.SRC0[i]), inst.Src0Sel)
			src1 := u.sdwaSrcSelect(uint32(sp.SRC1[i]), inst.Src1Sel)
			dst := src1 - src0
			dst = u.sdwaDstSelect(uint32(sp.DST[i]), dst,
				inst.DstSel, inst.DstUnused)
			sp.DST[i] = uint64(dst)
		}
	}
}

func (u *ALUImpl) runVADDCOCIU32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	newVCC := uint64(0)
	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			carry := (sp.VCC & (1 << i)) >> i

			if sp.SRC0[i] > math.MaxUint32-carry-sp.SRC1[i] {
				newVCC |= 1 << uint32(i)
			}

			sp.DST[i] = sp.SRC0[i] + sp.SRC1[i] + carry
			sp.DST[i] &= 0xffffffff
		}
		sp.VCC = newVCC
	} else {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := u.sdwaSrcSelect(uint32(sp.SRC0[i]), inst.Src0Sel)
			src1 := u.sdwaSrcSelect(uint32(sp.SRC1[i]), inst.Src1Sel)
			dst := src0 + src1 + uint32(sp.VCC)
			dst = u.sdwaDstSelect(uint32(sp.DST[i]), dst,
				inst.DstSel, inst.DstUnused)
			sp.DST[i] = uint64(dst)
			if src0 > math.MaxUint32-src1 {
				sp.VCC = 1
			} else {
				sp.VCC = 0
			}
		}
	}
}

func (u *ALUImpl) runVFMACF32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	var dst float32
	var src0 float32
	var src1 float32

	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			dst = asFloat32(uint32(sp.DST[i]))
			src0 = asFloat32(uint32(sp.SRC0[i]))
			src1 = asFloat32(uint32(sp.SRC1[i]))
			dst += src0 * src1
			sp.DST[i] = uint64(float32ToBits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode  %d not implemented \n", inst.Opcode)
	}
}

func (u *ALUImpl) runVORB32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := uint32(sp.SRC1[i])
			dst := src0 | src1
			sp.DST[i] = uint64(dst)
		}
	} else {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := u.sdwaSrcSelect(uint32(sp.SRC0[i]), inst.Src0Sel)
			src1 := u.sdwaSrcSelect(uint32(sp.SRC1[i]), inst.Src1Sel)
			dst := src0 | src1
			dst = u.sdwaDstSelect(uint32(sp.DST[i]), dst,
				inst.DstSel, inst.DstUnused)
			sp.DST[i] = uint64(dst)
		}
	}
}
func (u *ALUImpl) runVXORB32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	var i uint
	if !inst.IsSdwa {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := uint32(sp.SRC0[i])
			src1 := uint32(sp.SRC1[i])
			dst := src0 ^ src1
			sp.DST[i] = uint64(dst)
		}
	} else {
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}
			src0 := u.sdwaSrcSelect(uint32(sp.SRC0[i]), inst.Src0Sel)
			src1 := u.sdwaSrcSelect(uint32(sp.SRC1[i]), inst.Src1Sel)
			dst := src0 ^ src1
			dst = u.sdwaDstSelect(uint32(sp.DST[i]), dst,
				inst.DstSel, inst.DstUnused)
			sp.DST[i] = uint64(dst)
		}
	}
}
func (u *ALUImpl) runVMAXI32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := asInt32(uint32(sp.SRC0[i]))
			src1 := asInt32(uint32(sp.SRC1[i]))
			dst := src0
			if src1 > src0 {
				dst = src1
			}

			sp.DST[i] = uint64(int32ToBits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode %d not implemented \n", inst.Opcode)
	}
}
func (u *ALUImpl) runVMINI32(state InstEmuState) {
	sp := state.Scratchpad().AsVOP2()
	inst := state.Inst()
	if !inst.IsSdwa {
		var i uint
		for i = 0; i < 64; i++ {
			if !laneMasked(sp.EXEC, i) {
				continue
			}

			src0 := asInt32(uint32(sp.SRC0[i]))
			src1 := asInt32(uint32(sp.SRC1[i]))
			dst := src0
			if src1 < src0 {
				dst = src1
			}

			sp.DST[i] = uint64(int32ToBits(dst))
		}
	} else {
		log.Panicf("SDWA for VOP2 instruction opcode %d not implemented \n", inst.Opcode)
	}
}
