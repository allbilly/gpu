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

		state = new(mockInstState)
		state.scratchpad = make([]byte, 4096)
	})

	It("should run v_add_co_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3b
		state.inst.Opcode = 783

		sp := state.Scratchpad().AsVOP3B()
		sp.SRC0[0] = 0x00ff0000
		sp.SRC1[0] = 0x00ffff00
		sp.SRC0[1] = 0xffff0000
		sp.SRC1[1] = 0xf0ffff00
		sp.EXEC = 3

		alu.Run(state)

		Expect(sp.SDST).To(Equal(uint64(0x00000002)))
		Expect(sp.DST[0]).To(Equal(uint64(0x01feff00)))
		Expect(sp.DST[1]).To(Equal(uint64(0xf0feff00)))
	})

	It("should v_add_co_ci_u32", func() {
		state.inst = rdnainsts.NewInst()
		state.inst.FormatType = rdnainsts.VOP3b
		state.inst.Opcode = 296

		sp := state.Scratchpad().AsVOP3B()
		sp.SRC0[0] = math.MaxUint32 - 10
		sp.SRC1[0] = 10
		sp.SRC2[0] = 0x00000001
		sp.EXEC = 1

		alu.Run(state)

		Expect((sp.SDST)).To(Equal(uint64(0x00000001)))
		Expect(sp.DST[0]).To(Equal(uint64(0)))

	})

})
