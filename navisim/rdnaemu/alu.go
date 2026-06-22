package rdnaemu

import (
	"log"

	"gitlab.com/akita/navisim/rdnainsts"
)

//ALU does its jobs
type ALU interface {
	Run(state InstEmuState)

	SetLDS(lds []byte)
	LDS() []byte
}

// ALUImpl is where the instructions get executed.
type ALUImpl struct {
	storageAccessor *storageAccessor
	lds             []byte
}

// NewALU creates a new ALU with a storage as a dependency.
func NewALU(storageAccessor *storageAccessor) *ALUImpl {
	alu := new(ALUImpl)
	alu.storageAccessor = storageAccessor
	return alu
}

// SetLDS assigns the LDS storage to be used in the following instructions.
func (u *ALUImpl) SetLDS(lds []byte) {
	u.lds = lds
}

//LDS returns lds
func (u *ALUImpl) LDS() []byte {
	return u.lds
}

// Run executes the instruction in the scatchpad of the InstEmuState
//nolint:gocyclo
func (u *ALUImpl) Run(state InstEmuState) {
	inst := state.Inst()
	//fmt.Printf("%s\n", inst.String(nil))

	switch inst.FormatType {
	case rdnainsts.SOP1:
		u.runSOP1(state)
	case rdnainsts.SOP2:
		u.runSOP2(state)
	case rdnainsts.SOPC:
		u.runSOPC(state)
	case rdnainsts.SMEM:
		u.runSMEM(state)
	case rdnainsts.VOP1:
		u.runVOP1(state)
	case rdnainsts.VOP2:
		u.runVOP2(state)
	case rdnainsts.VOP3a:
		u.runVOP3A(state)
	case rdnainsts.VOP3b:
		u.runVOP3B(state)
	case rdnainsts.VOPC:
		u.runVOPC(state)
	case rdnainsts.FLAT:
		u.runFlat(state)
	case rdnainsts.SOPP:
		u.runSOPP(state)
	case rdnainsts.SOPK:
		u.runSOPK(state)
	case rdnainsts.DS:
		u.runDS(state)
	case rdnainsts.MUBUF:

	default:
		log.Panicf("Inst format %s is not supported", inst.Format.FormatName)
	}
}

func (u *ALUImpl) runSMEM(state InstEmuState) {
	inst := state.Inst()
	switch inst.Opcode {
	case 0:
		u.runSLOADDWORD(state)
	case 1:
		u.runSLOADDWORDX2(state)
	case 2:
		u.runSLOADDWORDX4(state)
	case 3:
		u.runSLOADDWORDX8(state)

	default:
		log.Panicf("Opcode %d for SMEM format is not implemented", inst.Opcode)
	}
}
func (u *ALUImpl) runSLOADDWORD(state InstEmuState) {
	sp := state.Scratchpad().AsSMEM()
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, sp.Base+sp.Offset, 4)

	sp.DST[0] = rdnainsts.BytesToUint32(buf)
}

func (u *ALUImpl) runSLOADDWORDX2(state InstEmuState) {
	sp := state.Scratchpad().AsSMEM()
	spRaw := state.Scratchpad()
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, sp.Base+sp.Offset, 8)
	copy(spRaw[32:40], buf)
}

func (u *ALUImpl) runSLOADDWORDX4(state InstEmuState) {
	sp := state.Scratchpad().AsSMEM()
	spRaw := state.Scratchpad()
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, sp.Base+sp.Offset, 16)
	copy(spRaw[32:48], buf)
}

func (u *ALUImpl) runSLOADDWORDX8(state InstEmuState) {
	sp := state.Scratchpad().AsSMEM()
	spRaw := state.Scratchpad()
	pid := state.PID()

	buf := u.storageAccessor.Read(pid, sp.Base+sp.Offset, 32)
	copy(spRaw[32:64], buf)
}

func (u *ALUImpl) sdwaSrcSelect(src uint32, sel rdnainsts.SDWASelect) uint32 {
	return src
}
func (u *ALUImpl) sdwaDstSelect(
	dstOld uint32,
	dstNew uint32,
	sel rdnainsts.SDWASelect,
	unused rdnainsts.SDWAUnused,
) uint32 {
	return dstNew
}
