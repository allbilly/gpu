package rdnaemu

import (
	"log"

	"gitlab.com/akita/navisim/rdnainsts"
)

//nolint:gocyclo
//nolint:funlen
func (u *ALUImpl) runFlat(state InstEmuState) {
	inst := state.Inst()
	//log.Printf("%s\n", inst.InstName)
	switch inst.Opcode {
	case 10:
		u.runFlatLoadUShort(state)
	case 12:
		u.runFlatLoadDWord(state)
	case 14:
		u.runFlatLoadDwordx4(state)
	case 28:
		u.runFlatStoreDWord(state)
	case 13:
		u.runFlatLoadDWordX2(state)
	case 30:
		u.runFlatStoreDWordx4(state)
	default:
		log.Panicf("Opcode %d for FLAT format is not implemented", inst.Opcode)
	}
}
func (u *ALUImpl) runFlatLoadUShort(state InstEmuState) {
	sp := state.Scratchpad().AsFlat()
	pid := state.PID()
	inst := state.Inst()

	for i := uint(0); i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		buf := u.storageAccessor.Read(pid, sp.ADDR[i]+uint64(inst.Offset.IntValue), uint64(4))
		buf[2] = 0
		buf[3] = 0
		sp.DST[i*4] = rdnainsts.BytesToUint32(buf)
	}
}

func (u *ALUImpl) runFlatLoadDWord(state InstEmuState) {
	sp := state.Scratchpad().AsFlat()
	pid := state.PID()
	inst := state.Inst()

	for i := uint(0); i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}
		buf := u.storageAccessor.Read(pid, sp.ADDR[i]+uint64(inst.Offset.IntValue), uint64(4))
		sp.DST[i*4] = rdnainsts.BytesToUint32(buf)
	}
}

func (u *ALUImpl) runFlatLoadDwordx4(state InstEmuState) {
	sp := state.Scratchpad().AsFlat()
	pid := state.PID()
	for i := uint(0); i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		buf := u.storageAccessor.Read(pid, sp.ADDR[i], uint64(16))

		sp.DST[i*4] = rdnainsts.BytesToUint32(buf[0:4])
		sp.DST[i*4+1] = rdnainsts.BytesToUint32(buf[4:8])
		sp.DST[i*4+2] = rdnainsts.BytesToUint32(buf[8:12])
		sp.DST[i*4+3] = rdnainsts.BytesToUint32(buf[12:16])
	}
}

func (u *ALUImpl) runFlatStoreDWord(state InstEmuState) {
	sp := state.Scratchpad().AsFlat()
	pid := state.PID()
	inst := state.Inst()

	for i := uint(0); i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		u.storageAccessor.Write(
			pid, sp.ADDR[i]+uint64(inst.Offset.IntValue), rdnainsts.Uint32ToBytes(sp.DATA[i*4]))
	}
}
func (u *ALUImpl) runFlatLoadDWordX2(state InstEmuState) {
	sp := state.Scratchpad().AsFlat()
	pid := state.PID()
	for i := uint(0); i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		buf := u.storageAccessor.Read(pid, sp.ADDR[i], uint64(8))
		sp.DST[i*4] = rdnainsts.BytesToUint32(buf[0:4])
		sp.DST[i*4+1] = rdnainsts.BytesToUint32(buf[4:8])
	}
}
func (u *ALUImpl) runFlatStoreDWordx4(state InstEmuState) {
	sp := state.Scratchpad().AsFlat()
	pid := state.PID()
	for i := uint(0); i < 64; i++ {
		if !laneMasked(sp.EXEC, i) {
			continue
		}

		buf := make([]byte, 16)
		copy(buf[0:4], rdnainsts.Uint32ToBytes(sp.DATA[i*4]))
		copy(buf[4:8], rdnainsts.Uint32ToBytes(sp.DATA[(i*4)+1]))
		copy(buf[8:12], rdnainsts.Uint32ToBytes(sp.DATA[(i*4)+2]))
		copy(buf[12:16], rdnainsts.Uint32ToBytes(sp.DATA[(i*4)+3]))

		u.storageAccessor.Write(pid, sp.ADDR[i], buf)
	}
}
