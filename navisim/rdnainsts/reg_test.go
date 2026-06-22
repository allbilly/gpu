package rdnainsts_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/rdnainsts"
)

var _ = Describe("Register", func() {
	It("should get correct v register", func() {
		Expect(rdnainsts.VReg(0)).To(BeIdenticalTo(rdnainsts.Regs[rdnainsts.V0]))
		Expect(rdnainsts.VReg(5)).To(BeIdenticalTo(rdnainsts.Regs[rdnainsts.V5]))
	})

	It("should get correct s register", func() {
		Expect(rdnainsts.SReg(0)).To(BeIdenticalTo(rdnainsts.Regs[rdnainsts.S0]))
		Expect(rdnainsts.SReg(5)).To(BeIdenticalTo(rdnainsts.Regs[rdnainsts.S5]))
	})

})
