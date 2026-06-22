package rdnaemu

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/rdnainsts"
)

var _ = Describe("Wavefront", func() {
	var (
		wf *Wavefront
	)

	BeforeEach(func() {
		wf = NewWavefront(nil)
	})

	It("should write vcc_lo", func() {
		wf.VCC = 0xf

		wf.WriteReg(
			rdnainsts.Regs[rdnainsts.VCCLO],
			1,
			0,
			[]byte{0, 0, 0, 0},
		)

		Expect(wf.VCC).To(Equal(uint64(0)))
	})

	It("should write exec_lo", func() {
		wf.Exec = 0xf

		wf.WriteReg(
			rdnainsts.Regs[rdnainsts.EXECLO],
			1,
			0,
			[]byte{0, 0, 0, 0},
		)

		Expect(wf.Exec).To(Equal(uint64(0)))
	})
})
