package cu

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/kernels"
	"gitlab.com/akita/navisim/protocol"
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/navisim/timing/wavefront"
)

var _ = Describe("WfDispatcher", func() {
	var (
		cu           *ComputeUnit
		wfDispatcher *WfDispatcherImpl
	)

	BeforeEach(func() {
		cu = NewComputeUnit("cu", nil)
		cu.Freq = 1

		for i := 0; i < 4; i++ {
			sRegFile := NewSimpleRegisterFile(uint64(2560*4), 0)
			vRegFile := NewSimpleRegisterFile(uint64(32768*4), 1024)
			cu.VRegFiles = append(cu.VRegFiles, vRegFile)
			cu.SRegFiles = append(cu.SRegFiles, sRegFile)
		}

		wfDispatcher = NewWfDispatcher(cu)
	})

	It("should dispatch wavefront", func() {
		rawWf := kernels.NewWavefront(32)
		rawWG := kernels.NewWorkGroup()
		rawWf.WG = rawWG
		rawWG.SizeX = 256
		rawWG.SizeY = 1
		rawWG.SizeZ = 1
		wfDispatchInfo := protocol.WfDispatchLocation{
			Wavefront:  rawWf,
			SIMDID:     1,
			VGPROffset: 16,
			SGPROffset: 8,
			LDSOffset:  512,
		}

		co := rdnainsts.NewHsaCo()
		co.KernelCodeEntryByteOffset = 256
		packet := new(kernels.HsaKernelDispatchPacket)
		packet.KernelObject = 65536

		wf := wavefront.NewWavefront(rawWf)
		wg := wavefront.NewWorkGroup(rawWG, nil)
		wf.WG = wg
		wf.CodeObject = co
		wf.Packet = packet
		//req := navisim.NewDispatchWfReq(nil, cu.ToACE, 10, nil)
		wfDispatcher.DispatchWf(10, wf, wfDispatchInfo)

		//Expect(len(engine.ScheduledEvent)).To(Equal(1))
		Expect(wf.SIMDID).To(Equal(1))
		Expect(wf.VRegOffset).To(Equal(16))
		Expect(wf.SRegOffset).To(Equal(8))
		Expect(wf.LDSOffset).To(Equal(512))
		Expect(wf.PC).To(Equal(uint64(65536 + 256)))
	})
})
