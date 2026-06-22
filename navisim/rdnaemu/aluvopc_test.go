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

	It("should run v_cmp_lt_f32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 1

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0xF
		sp.SRC0[0] = uint64(math.Float32bits(-1.0))
		sp.SRC0[1] = uint64(math.Float32bits(-1.0))
		sp.SRC0[2] = uint64(math.Float32bits(-1.0))
		sp.SRC0[3] = uint64(math.Float32bits(-1.0))
		sp.SRC1[0] = uint64(math.Float32bits(-1.0))
		sp.SRC1[1] = uint64(math.Float32bits(-2.0))
		sp.SRC1[2] = uint64(math.Float32bits(1.0))
		sp.SRC1[3] = uint64(math.Float32bits(2.0))

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0b1100)))
	})

	It("should run v_cmp_gt_i32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 132

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0xF
		sp.SRC0[0] = 1
		sp.SRC0[1] = uint64(int32ToBits(-1))
		sp.SRC0[2] = 1
		sp.SRC0[3] = 1
		sp.SRC1[0] = 1
		sp.SRC1[1] = uint64(int32ToBits(-2))
		sp.SRC1[2] = 0
		sp.SRC1[3] = 2

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x6)))
	})

	It("should run v_cmp_le_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 195

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0xffffffffffffffff
		sp.SRC0[0] = 1
		sp.SRC0[1] = 1
		sp.SRC0[2] = 1
		sp.SRC1[0] = 1
		sp.SRC1[1] = 2
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0xfffffffffffffffb)))
	})

	It("should run v_cmp_gt_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 196

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = 1
		sp.SRC1[0] = 1
		sp.SRC0[1] = 1
		sp.SRC1[1] = 2
		sp.SRC0[2] = 1
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x4)))
	})

	It("should run v_cmp_lt_i32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 129

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = uint64(int32ToBits(math.MinInt32))
		sp.SRC1[0] = 1
		sp.SRC0[1] = uint64(int32ToBits(math.MaxInt32))
		sp.SRC1[1] = 2
		sp.SRC0[2] = 0
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x1)))
	})
	
	It("should run v_cmp_ne_u64", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 229

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = 1
		sp.SRC1[0] = 1
		sp.SRC0[1] = 1
		sp.SRC1[1] = 2
		sp.SRC0[2] = 1
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x6)))
	})

	It("should run v_cmp_gt_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 198

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = 1
		sp.SRC1[0] = 1
		sp.SRC0[1] = 1
		sp.SRC1[1] = 2
		sp.SRC0[2] = 1
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x5)))
	})

	It("should run v_cmp_le_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 195

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0xffffffffffffffff
		sp.SRC0[0] = 1
		sp.SRC0[1] = 1
		sp.SRC0[2] = 1
		sp.SRC1[0] = 1
		sp.SRC1[1] = 2
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0xfffffffffffffffb)))
	})

	It("should run v_cmp_ne_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 197

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0xffffffffffffffff
		sp.SRC0[0] = 1
		sp.SRC1[0] = 1
		sp.SRC0[1] = 0
		sp.SRC1[1] = 2

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x0000000000000002)))
	})

	It("should run v_cmp_eq_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 194

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = 1
		sp.SRC0[1] = 1
		sp.SRC0[2] = 1
		sp.SRC1[0] = 1
		sp.SRC1[1] = 2
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x1)))
	})

	It("should run v_cmp_le_i32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 131

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = 0
		sp.SRC1[0] = 1
		sp.SRC0[1] = 2
		sp.SRC1[1] = 2
		sp.SRC0[2] = 1
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x3)))
	})
	
	It("should run v_cmp_lt_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOPC
		state.inst.Opcode = 193

		sp := state.Scratchpad().AsVOPC()
		sp.EXEC = 0x7
		sp.SRC0[0] = 0
		sp.SRC1[0] = 1
		sp.SRC0[1] = 2
		sp.SRC1[1] = 2
		sp.SRC0[2] = 1
		sp.SRC1[2] = 0

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x1)))
	})

})
