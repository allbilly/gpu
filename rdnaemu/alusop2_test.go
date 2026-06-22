package rdnaemu

import (
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

	It("should run S_ADD_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 0

		sp := state.scratchpad.AsSOP2()
		sp.SRC0 = 1<<32 - 9
		sp.SRC1 = 10

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(1)))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ADD_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 0

		sp := state.scratchpad.AsSOP2()
		sp.SRC0 = 1<<32 - 9
		sp.SRC1 = 10

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(1)))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_SUB_I32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.scratchpad.AsSOP2()
		sp.SRC0 = 15
		sp.SRC1 = 3

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(12)))
		Expect(sp.SCC).To(Equal(byte(0)))
	})

	It("should run S_MIN_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 7

		sp := state.scratchpad.AsSOP2()
		sp.SRC0 = 1
		sp.SRC1 = 2

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(1)))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ADD_I32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 2

		sp := state.scratchpad.AsSOP2()
		sp.SRC0 = 1<<31 - 9
		sp.SRC1 = 10

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(2147483649)))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ADDC_U32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 4

		sp := state.scratchpad.AsSOP2()
		sp.SRC0 = 1<<32 - 9
		sp.SRC1 = 10
		sp.SCC = 1

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(2)))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_LSHL_B64", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 31

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0xAAAA_AAAA_AAAA_AAAA
		sp.SRC1 = 2

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(0xAAAA_AAAA_AAAA_AAA8)))
	})

	It("should run S_ASHR_I32 (Negative)", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 34

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = int64ToBits(-128)
		sp.SRC1 = 2

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(int32ToBits(-32))))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ASHR_I32 (Positive)", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 34

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = int64ToBits(128)
		sp.SRC1 = 2

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(int32ToBits(32))))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})
	It("should run S_XOR_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 18

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0xf0
		sp.SRC1 = 0xff

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(0x0f)))
		Expect(sp.SCC).To(Equal(byte(1)))
	})
	It("should run S_SUB_I32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 10
		sp.SRC1 = 6

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(4)))
		Expect(sp.SCC).To(Equal(byte(0)))
	})

	It("should run S_SUB_I32, when input is negative", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = uint64((int32ToBits(-6)))
		sp.SRC1 = 15

		alu.Run(state)

		Expect(asInt32(uint32(sp.DST))).To(Equal(int32(-21)))
		Expect(sp.SCC).To(Equal(byte(0)))
	})

	It("should run S_SUB_I32, when overflow and src1 is positive", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0x7ffffffe
		sp.SRC1 = 0xfffffffc

		alu.Run(state)

		Expect(sp.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUB_I32, when overflow and src1 is negtive", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0x80000001
		sp.SRC1 = 10

		alu.Run(state)

		Expect(sp.SCC).To(Equal(byte(1)))
	})

	It("should run S_ASHR_I32 (Negative)", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 34

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = int64ToBits(-128)
		sp.SRC1 = 2

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(int32ToBits(-32))))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_ASHR_I32 (Positive)", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 34

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = int64ToBits(128)
		sp.SRC1 = 2

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(int32ToBits(32))))
		Expect(sp.SCC).To(Equal(uint8(1)))
	})

	It("should run S_XOR_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 18

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0xf0
		sp.SRC1 = 0xff

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(0x0f)))
		Expect(sp.SCC).To(Equal(byte(1)))
	})
	It("should run S_OR_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 16

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0xf0
		sp.SRC1 = 0xff

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(0xff)))
		Expect(sp.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUB_I32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 10
		sp.SRC1 = 6

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(4)))
		Expect(sp.SCC).To(Equal(byte(0)))
	})

	It("should run S_SUB_I32, when input is negative", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = uint64((int32ToBits(-6)))
		sp.SRC1 = 15

		alu.Run(state)

		Expect(asInt32(uint32(sp.DST))).To(Equal(int32(-21)))
		Expect(sp.SCC).To(Equal(byte(0)))
	})

	It("should run S_SUB_I32, when overflow and src1 is positive", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0x7ffffffe
		sp.SRC1 = 0xfffffffc

		alu.Run(state)

		Expect(sp.SCC).To(Equal(byte(1)))
	})

	It("should run S_SUB_I32, when overflow and src1 is negtive", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 3

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0x80000001
		sp.SRC1 = 10

		alu.Run(state)

		Expect(sp.SCC).To(Equal(byte(1)))
	})
	It("should run S_ANDN2_B64", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 20

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0xab
		sp.SRC1 = 0x0f

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(0xa0)))
		Expect(sp.SCC).To(Equal(byte(1)))
	})

	It("should run S_LSHR_B32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOP2
		state.inst.Opcode = 32

		sp := state.Scratchpad().AsSOP2()
		sp.SRC0 = 0x20
		sp.SRC1 = 0x64

		alu.Run(state)

		Expect(sp.DST).To(Equal(uint64(0x02)))
		Expect(sp.SCC).To(Equal(byte(1)))
	})
})
