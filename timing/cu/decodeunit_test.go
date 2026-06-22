package cu

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/timing/wavefront"
)

var _ = Describe("DecodeUnit", func() {
	var (
		cu        *ComputeUnit
		du        *DecodeUnit
		execUnit *mockCUComponent
	)

	BeforeEach(func() {
		cu = NewComputeUnit("cu", nil)
		du = NewDecodeUnit(cu)
		execUnit = new(mockCUComponent)
		execUnit.canAccept = true
		du.ExecUnit = execUnit
	})

	It("should tell if it cannot accept wave", func() {
		du.toDecode = new(wavefront.Wavefront)
		Expect(du.CanAcceptWave()).To(BeFalse())
	})

	It("should tell if it can accept wave", func() {
		du.toDecode = nil
		Expect(du.CanAcceptWave()).To(BeTrue())
	})

	It("should accept wave", func() {
		wave := new(wavefront.Wavefront)
		du.toDecode = nil
		du.AcceptWave(wave, 10)
		Expect(du.toDecode).To(BeIdenticalTo(wave))
	})

	It("should return error if the decoder is busy", func() {
		wave := new(wavefront.Wavefront)
		wave2 := new(wavefront.Wavefront)
		du.toDecode = wave

		Expect(func() { du.AcceptWave(wave2, 10) }).Should(Panic())
		Expect(du.toDecode).To(BeIdenticalTo(wave))
	})

	It("should deliver the wave to the execution unit", func() {
		wave := new(wavefront.Wavefront)
		du.toDecode = wave

		du.Run(10)

		Expect(len(execUnit.acceptedWave)).To(Equal(1))
		Expect(du.toDecode).To(BeNil())
	})

	It("should not deliver to the execution unit, if busy", func() {
		wave := new(wavefront.Wavefront)
		du.toDecode = wave
		execUnit.canAccept = false

		du.Run(10)

		Expect(len(execUnit.acceptedWave)).To(Equal(0))
	})
	It("should flush the decode unit", func() {
		wave := new(wavefront.Wavefront)
		du.toDecode = wave

		du.Flush()

		Expect(du.toDecode).To(BeNil())
	})

})
