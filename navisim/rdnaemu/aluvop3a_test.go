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

	It("should run v_cmp_gt_f32 VOP3a", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 4

		sp := state.Scratchpad().AsVOP3A()
		sp.EXEC = 0x7
		sp.SRC0[0] = uint64(math.Float32bits(-1.2))
		sp.SRC1[0] = uint64(math.Float32bits(-1.2))
		sp.SRC0[1] = uint64(math.Float32bits(-2.5))
		sp.SRC1[1] = uint64(math.Float32bits(0.0))
		sp.SRC0[2] = uint64(math.Float32bits(1.5))
		sp.SRC1[2] = uint64(math.Float32bits(0.0))

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(0x4)))
	})

	It("should run v_cmp_ge_i32 VOP3a", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 134

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = uint64(int32ToBits(0))
		sp.SRC1[0] = uint64(int32ToBits(math.MinInt32))
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal((uint64(1))))
	})

	It("should run v_cmp_gt_i32 VOP3a", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 132

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = uint64(int32ToBits(0))
		sp.SRC1[0] = uint64(int32ToBits(math.MinInt32))
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal((uint64(1))))
	})

	It("should run V_CNDMASK_B32 VOP3a", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 257

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = 1
		sp.SRC1[0] = 2
		sp.SRC0[1] = 1
		sp.SRC1[1] = 2
		sp.SRC2[0] = 1
		sp.EXEC = 3

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(2)))
		Expect(sp.DST[1]).To(Equal(uint64(1)))
	})

	It("should run V_MAD_U32_U24 VOP3a", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 323

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = uint64(int32ToBits(10))
		sp.SRC1[0] = uint64(int32ToBits(20))
		sp.SRC2[0] = uint64(int32ToBits(50))
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0] & 0xffffffff).To(Equal(uint64(250)))
	})

	It("should run V_MUL_LO_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 361

		sp := state.Scratchpad().AsVOP3A()
		for i := 0; i < 64; i++ {
			sp.SRC0[i] = uint64(i)
			sp.SRC1[i] = uint64(2)
		}
		sp.EXEC = 0xffffffffffffffff

		alu.Run(state)

		for i := 0; i < 64; i++ {
			Expect(sp.DST[i]).To(Equal(uint64(i * 2)))
		}
	})

	It("should run V_MUL_HI_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 362

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = uint64(0x80000000)
		sp.SRC1[0] = uint64(2)
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(1)))

	})

	It("should run V_LSHL_REV_B64", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 767

		sp := state.Scratchpad().AsVOP3A()

		sp.SRC1[0] = uint64(0x0000000000010000)
		sp.SRC0[0] = uint64(3)
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(0x0000000000080000)))

	})

	It("should run V_MUL_LO_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 361

		sp := state.Scratchpad().AsVOP3A()
		for i := 0; i < 64; i++ {
			sp.SRC0[i] = uint64(i)
			sp.SRC1[i] = uint64(2)
		}
		sp.EXEC = 0xffffffffffffffff

		alu.Run(state)

		for i := 0; i < 64; i++ {
			Expect(sp.DST[i]).To(Equal(uint64(i * 2)))
		}
	})

	It("should run V_MUL_HI_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 362

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = uint64(0x80000000)
		sp.SRC1[0] = uint64(2)
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(1)))
	})

	It("should run V_LSHL_ADD_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 838

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = uint64(2)
		sp.SRC1[0] = uint64(2)
		sp.SRC2[0] = uint64(1)
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(9)))
	})

	It("should run V_MUL_U32_U24", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 267

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = 10
		sp.SRC1[0] = 20
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(200)))
	})

	It("should run V_MAD_U32_U24", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 323

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = 10
		sp.SRC1[0] = 20
		sp.SRC2[0] = 50
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(250)))
	})

	It("should run v_cmp_eq_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 194

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = 10
		sp.SRC1[0] = 20
		sp.SRC2[0] = 50
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(0)))
	})

	It("should run V_ADD3_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3a
		state.inst.Opcode = 877

		sp := state.Scratchpad().AsVOP3A()
		sp.SRC0[0] = 5
		sp.SRC1[0] = 2
		sp.SRC2[0] = 1
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(8)))
	})
})
