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

	It("should run V_SUB_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 4

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(math.Float32bits(2.0))
		sp.SRC1[0] = uint64(math.Float32bits(3.1))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(
			BeNumerically("~", -1.1, 1e-4))
	})

	It("should run V_SUBREV_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 5

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(math.Float32bits(2.0))
		sp.SRC1[0] = uint64(math.Float32bits(3.1))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(math.Float32frombits(uint32(sp.DST[0]))).To(
			BeNumerically("~", 1.1, 1e-4))
	})

	It("should run V_MUL_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 8

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(math.Float32bits(2.0))
		sp.SRC1[0] = uint64(math.Float32bits(4.2))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(math.Float32bits(float32(8.4)))))
	})

	It("should run V_ASHRREV_I32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 24

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 97
		sp.SRC1[0] = uint64(int32ToBits(-64))
		sp.EXEC = 1

		alu.Run(state)
		Expect(asInt32(uint32(sp.DST[0]))).To(Equal(int32(-32)))

	})

	It("should run V_LSHRREV_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 26

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 0x64
		sp.SRC1[0] = 0x02
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(uint32(sp.DST[0])).To(Equal(uint32(0x20)))
	})

	It("should run V_CNDMASK_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 1

		sp := state.Scratchpad().AsVOP2()
		sp.VCC = 1
		sp.SRC0[0] = 1
		sp.SRC0[1] = 2
		sp.SRC1[0] = 3
		sp.SRC1[1] = 4
		sp.EXEC = 3

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(3)))
		Expect(sp.DST[1]).To(Equal(uint64(2)))
		Expect(sp.DST[1]).To(Equal(uint64(2)))
	})

	It("should run V_FMAC_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 43

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(float32ToBits(4))
		sp.SRC1[0] = uint64(float32ToBits(16))
		sp.DST[0] = uint64(float32ToBits(1024))
		sp.EXEC = 1

		alu.Run(state)

		Expect(asFloat32(uint32(sp.DST[0]))).To(Equal(float32(1024.0 + 16.0*4.0)))
	})
	It("should run v_max_f32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 16

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(float32ToBits(4))
		sp.SRC1[0] = uint64(float32ToBits(16))
		sp.EXEC = 1

		alu.Run(state)

		Expect(asFloat32(uint32(sp.DST[0]))).To(Equal(float32(16)))
	})
	It("should run v_add_co_ci", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 40

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 0xf0000000
		sp.SRC1[0] = 0x0000f000
		sp.SRC0[1] = 0xf0000000
		sp.SRC1[1] = 0x1000f000
		sp.VCC = 0x00000001
		sp.EXEC = 3

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0x00000002)))
		Expect(sp.DST[0]).To(Equal(uint64(0xf000f001)))
		Expect(sp.DST[1]).To(Equal(uint64(0x0000f000)))
	})

	It("should run v_add_nc_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 37

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = math.MaxUint32 - 10
		sp.SRC1[0] = 11
		sp.VCC = 0
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.VCC).To(Equal(uint64(0)))
		Expect(sp.DST[0]).To(Equal(uint64(0)))
	})

	It("should run v_or_b32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 28

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 0x0000ffff
		sp.SRC1[0] = 0xff000000
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(0xff00ffff)))
	})
	It("should run V_MUL_F32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 8

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(math.Float32bits(2.0))
		sp.SRC1[0] = uint64(math.Float32bits(3.1))
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(math.Float32bits(float32(6.2)))))
	})

	It("should run V_ASHRREV_I32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 24

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 97
		sp.SRC1[0] = uint64(int32ToBits(-64))
		sp.EXEC = 1

		alu.Run(state)
		Expect(asInt32(uint32(sp.DST[0]))).To(Equal(int32(-32)))

	})

	It("should run V_MUL_U32_U24", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 11

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 2
		sp.SRC1[0] = 0x1000001
		sp.EXEC = 0x1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(2)))
	})
	It("should run v_max_i32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 18

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(int32ToBits(math.MinInt32))
		sp.SRC1[0] = uint64(int32ToBits(math.MinInt32 + 1))
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(int32ToBits(math.MinInt32 + 1))))
	})
	It("should run v_min_i32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 17

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(int32ToBits(math.MinInt32))
		sp.SRC1[0] = 0
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal(uint64(int32ToBits(math.MinInt32))))
	})
	It("should run v_add_f32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(float32ToBits(0.25))
		sp.SRC1[0] = uint64(float32ToBits(4.25))
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal((uint64(float32ToBits(4.5)))))
	})
	It("should run v_min_f32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 15

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = uint64(float32ToBits(0.25))
		sp.SRC1[0] = uint64(float32ToBits(4.25))
		sp.EXEC = 1

		alu.Run(state)

		Expect(sp.DST[0]).To(Equal((uint64(float32ToBits(0.25)))))
	})
	It("should run v_lshlrev_b32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 26

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 0x64
		sp.SRC1[0] = 0x02
		sp.EXEC = 1

		alu.Run(state)

		Expect(uint32(sp.DST[0])).To(Equal(uint32(0x20)))
	})
	It("should run V_XOR_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 29

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 2 // 10
		sp.SRC1[0] = 3 // 11
		sp.EXEC = 1

		alu.Run(state)

		Expect(uint32(sp.DST[0])).To(Equal(uint32(1)))
	})
	It("should run v_sub_nc_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 38

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 11
		sp.SRC1[0] = 10
		sp.EXEC = 1
		alu.Run(state)

		Expect(uint32(sp.DST[0])).To(Equal(uint32(1)))

	})
	It("should run v_sub_rev_nc_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP2
		state.inst.Opcode = 39

		sp := state.Scratchpad().AsVOP2()
		sp.SRC0[0] = 1
		sp.SRC1[0] = 10
		sp.EXEC = 1
		alu.Run(state)

		Expect(uint32(sp.DST[0])).To(Equal(uint32(9)))

	})
})
