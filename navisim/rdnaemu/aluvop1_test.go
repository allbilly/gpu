package rdnaemu

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/rdnainsts"
)

var _ = Describe("ALU", func() {

	var (
		alu   *ALUImpl
		state *mockInstState
	)

	BeforeEach(func() {
		alu = NewALU(nil)

		state = new(mockInstState)
		state.scratchpad = make([]byte, 4096)
	})

	It("should run V_CVT_F32_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 6

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = 1
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(1.0)))
	})
	It("should run V_RCP_IFLAG_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 43

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(1.0))
		sp.SRC0[1] = uint64(math.Float32bits(2.0))
		sp.EXEC = 0x3

		alu.Run(state)

		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(Equal(float32(1.0)))
		Expect(math.Float32frombits(uint32(sp.DST[1]))).To(Equal(float32(0.5)))
	})

	It("should run V_CVT_U32_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 7

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(1.0))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(1)))
	})

	It("should run V_CVT_U32_F32, when input is nan", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 7

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(float32(math.NaN())))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(0)))
	})

	It("should run V_CVT_U32_F32, when the input is negative", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 7

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(-1.0))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(0)))
	})

	It("should run V_CVT_U32_F32, when the input is very large", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 7

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(float32(math.MaxUint32 + 1)))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(math.MaxUint32)))
	})


	It("should run V_CVT_U32_F32, when the input is very large", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 7

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(float32(math.MaxUint32 + 1)))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(math.MaxUint32)))
	})


	It("should run V_CVT_I32_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP1
		state.inst.Opcode = 8

		sp := state.Scratchpad().AsVOP1()
		sp.SRC0[0] = uint64(math.Float32bits(float32(math.MaxUint32 + 1)))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(math.MaxInt32)))
	})

})
