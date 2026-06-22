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

	It("should run S_CMP_EQ_I32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 0

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = uint64(int32ToBits(5))
		layout.SRC1 = uint64(int32ToBits(5))

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_I32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 1

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = uint64(int32ToBits(3))
		layout.SRC1 = uint64(int32ToBits(5))

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_GT_I32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 2

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = uint64(int32ToBits(7))
		layout.SRC1 = uint64(int32ToBits(5))

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_GE_I32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 3

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = uint64(int32ToBits(7))
		layout.SRC1 = uint64(int32ToBits(5))

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LT_I32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 4

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = uint64(int32ToBits(5))
		layout.SRC1 = uint64(int32ToBits(7))

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LE_I32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 5

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = uint64(int32ToBits(5))
		layout.SRC1 = uint64(int32ToBits(7))

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_EQ_U32 when input is not equal", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 6

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = 1
		layout.SRC1 = 2

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_EQ_U32 when input is equal", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 6

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = 1
		layout.SRC1 = 1

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_U32 when condition holds", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 7

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = 1
		layout.SRC1 = 2

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

	It("should run S_CMP_LG_U32 when condition does not hold", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 7

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = 1
		layout.SRC1 = 1

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(0)))
	})

	It("should run S_CMP_GE_U32 when input is greater", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.SOPC
		state.inst.Opcode = 9

		layout := state.Scratchpad().AsSOPC()
		layout.SRC0 = 10
		layout.SRC1 = 2

		alu.Run(state)

		Expect(layout.SCC).To(Equal(byte(1)))
	})

})
