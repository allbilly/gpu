package resource

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/kernels"
	"gitlab.com/akita/navisim/rdnainsts"
)

func assertAllResourcesFree(r *CUResourceImpl) {
	Expect(r.wfPoolFreeCount[0]).To(Equal(20))
	Expect(r.wfPoolFreeCount[1]).To(Equal(20))
	Expect(r.wfPoolFreeCount[2]).To(Equal(20))
	Expect(r.wfPoolFreeCount[3]).To(Equal(20))
	Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(512))
	Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(160))
	Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(160))
	Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(160))
	Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(160))
	Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(256))
	Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(256))
	Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(256))
	Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(256))
}

var _ = Describe("cuResource", func() {
	var (
		r      *CUResourceImpl
		wg     *kernels.WorkGroup
		co     *rdnainsts.HsaCo
		wfSize int
	)

	BeforeEach(func() {
		r = &CUResourceImpl{
			wfPoolFreeCount: []int{20, 20, 20, 20},
			sregCounts:      []int{2560, 2560, 2560, 2560},
			sregGranularity: 16,
			sregMasks: []resourceMask{
				newResourceMask(2560 / 16),
				newResourceMask(2560 / 16),
				newResourceMask(2560 / 16),
				newResourceMask(2560 / 16),
			},
			vregCounts:      []int{1024, 1024, 1024, 1024},
			vregGranularity: 4,
			vregMasks: []resourceMask{
				newResourceMask(1024 / 4),
				newResourceMask(1024 / 4),
				newResourceMask(1024 / 4),
				newResourceMask(1024 / 4),
			},
			ldsByteSize:    128 * 1024,
			ldsGranularity: 256,
			ldsMask:        newResourceMask(128 * 1024 / 256),
			reservedWGs:    make(map[*kernels.WorkGroup][]WfLocation),
		}

		wg = kernels.NewWorkGroup()
		wfSize = 32
		for i := 0; i < 20; i++ {
			wf := kernels.NewWavefront(wfSize)
			wg.Wavefronts = append(wg.Wavefronts, wf)
		}

		co = rdnainsts.NewHsaCo()
		wg.CodeObject = co
	})

	It("should send NACK if too many Wavefronts", func() {
		// Each SIMD is running 18 wf in each SIMD. 2 more wfs can handle.
		for i := 0; i < 4; i++ {
			r.wfPoolFreeCount[i] = 2
		}

		info, ok := r.ReserveResourceForWG(wg)

		Expect(info).To(BeEmpty())
		Expect(ok).To(BeFalse())
		Expect(r.wfPoolFreeCount[0]).To(Equal(2))
		Expect(r.wfPoolFreeCount[1]).To(Equal(2))
		Expect(r.wfPoolFreeCount[2]).To(Equal(2))
		Expect(r.wfPoolFreeCount[3]).To(Equal(2))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(512))
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(256))
	})

	It("should send NACK to the dispatcher if too many SReg", func() {
		// 160 groups in total, 150 groups occupied.
		// 10 groups are free -> 160 registers available
		r.sregMasks[0].setStatus(0, 150, allocStatusReserved)
		r.sregMasks[1].setStatus(0, 150, allocStatusReserved)
		r.sregMasks[2].setStatus(0, 150, allocStatusReserved)
		r.sregMasks[3].setStatus(0, 150, allocStatusReserved)

		// 20 Wfs, 64 SGPRs per wf. That is 1280 in total
		co.WFSgprCount = 64

		_, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeFalse())
		Expect(r.wfPoolFreeCount[0]).To(Equal(20))
		Expect(r.wfPoolFreeCount[1]).To(Equal(20))
		Expect(r.wfPoolFreeCount[2]).To(Equal(20))
		Expect(r.wfPoolFreeCount[3]).To(Equal(20))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(512))
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(10))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(10))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(10))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(10))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(256))

	})

	It("should send NACK to the dispatcher if too large LDS", func() {
		// 496 units occupied, 16 units left -> 4096 Bytes available
		r.ldsMask.setStatus(0, 496, allocStatusReserved)

		co.WGGroupSegmentByteSize = 8192

		_, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeFalse())
		Expect(r.wfPoolFreeCount[0]).To(Equal(20))
		Expect(r.wfPoolFreeCount[1]).To(Equal(20))
		Expect(r.wfPoolFreeCount[2]).To(Equal(20))
		Expect(r.wfPoolFreeCount[3]).To(Equal(20))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(16))
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(256))
	})

	It("should send NACK if too many VGPRs", func() {
		// 64 units occupied, 4 units available, 4 * 4 = 16 units
		r.vregMasks[0].setStatus(0, 252, allocStatusReserved)
		r.vregMasks[1].setStatus(0, 252, allocStatusReserved)
		r.vregMasks[2].setStatus(0, 252, allocStatusReserved)
		r.vregMasks[3].setStatus(0, 252, allocStatusReserved)

		co.WFSgprCount = 20
		co.WGGroupSegmentByteSize = 256
		co.WIVgprCount = 20

		_, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeFalse())
		Expect(r.wfPoolFreeCount[0]).To(Equal(20))
		Expect(r.wfPoolFreeCount[1]).To(Equal(20))
		Expect(r.wfPoolFreeCount[2]).To(Equal(20))
		Expect(r.wfPoolFreeCount[3]).To(Equal(20))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(512))
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(4))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(4))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(4))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(4))
	})

	It("should send NACK if not all Wavefront can fit the VGPRs requirement", func() {
		// SIMD 0 and 1 do not have enouth VGPRs
		r.vregMasks[0].setStatus(0, 252, allocStatusReserved)
		r.vregMasks[1].setStatus(0, 252, allocStatusReserved)
		r.wfPoolFreeCount[2] = 2
		r.wfPoolFreeCount[3] = 2

		co.WIVgprCount = 102

		_, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeFalse())
		Expect(r.wfPoolFreeCount[0]).To(Equal(20))
		Expect(r.wfPoolFreeCount[1]).To(Equal(20))
		Expect(r.wfPoolFreeCount[2]).To(Equal(2))
		Expect(r.wfPoolFreeCount[3]).To(Equal(2))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(512))
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(160))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(4))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(4))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(256))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(256))
	})

	It("should reserve resources and send ACK back if all requirement satisfy", func() {
		co.WIVgprCount = 20
		co.WFSgprCount = 16
		co.WGGroupSegmentByteSize = 1024

		locations, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeTrue())
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[0].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[1].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[2].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[3].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(508))
		Expect(r.ldsMask.statusCount(allocStatusReserved)).To(Equal(4))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[0].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[1].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[2].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[3].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.wfPoolFreeCount[0]).To(Equal(15))
		Expect(r.wfPoolFreeCount[1]).To(Equal(15))
		Expect(r.wfPoolFreeCount[2]).To(Equal(15))
		Expect(r.wfPoolFreeCount[3]).To(Equal(15))

		for i := 0; i < len(wg.Wavefronts); i++ {
			Expect(locations[i].SIMDID).To(Equal(i % 4))
			Expect(locations[i].SGPROffset).To(Equal((i / 4) * 64))
			Expect(locations[i].LDSOffset).To(Equal(0))
			Expect(locations[i].VGPROffset).To(Equal((i / 4) * 20 * 4))
		}
	})

	It("should reserve resources if resources are not aligned with granularity", func() {
		co.WIVgprCount = 18
		co.WFSgprCount = 14
		co.WGGroupSegmentByteSize = 900

		locations, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeTrue())
		Expect(r.sregMasks[0].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[0].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.sregMasks[1].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[1].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.sregMasks[2].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[2].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.sregMasks[3].statusCount(allocStatusFree)).To(Equal(155))
		Expect(r.sregMasks[3].statusCount(allocStatusReserved)).To(Equal(5))
		Expect(r.ldsMask.statusCount(allocStatusFree)).To(Equal(508))
		Expect(r.ldsMask.statusCount(allocStatusReserved)).To(Equal(4))
		Expect(r.vregMasks[0].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[0].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.vregMasks[1].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[1].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.vregMasks[2].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[2].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.vregMasks[3].statusCount(allocStatusFree)).To(Equal(231))
		Expect(r.vregMasks[3].statusCount(allocStatusReserved)).To(Equal(25))
		Expect(r.wfPoolFreeCount[0]).To(Equal(15))
		Expect(r.wfPoolFreeCount[1]).To(Equal(15))
		Expect(r.wfPoolFreeCount[2]).To(Equal(15))
		Expect(r.wfPoolFreeCount[3]).To(Equal(15))

		for i := 0; i < len(wg.Wavefronts); i++ {
			Expect(locations[i].SIMDID).To(Equal(i % 4))
			Expect(locations[i].SGPROffset).To(Equal((i / 4) * 64))
			Expect(locations[i].LDSOffset).To(Equal(0))
			Expect(locations[i].VGPROffset).To(Equal((i / 4) * 20 * 4))
		}
	})

	It("should support non-standard CU size", func() {
		r.wfPoolFreeCount = []int{20, 20, 10, 10, 10}
		r.vregCounts = []int{1024, 1024, 1024, 1024, 1024}
		r.sregCounts = []int{1024, 1024, 1024, 1024, 1024}
		r.vregMasks = []resourceMask{
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
		}
		r.sregMasks = []resourceMask{
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
			newResourceMask(1024 / 4),
		}

		co.WIVgprCount = 20
		co.WFSgprCount = 16

		_, ok := r.ReserveResourceForWG(wg)

		Expect(ok).To(BeTrue())
		Expect(r.wfPoolFreeCount[0]).To(Equal(16))
		Expect(r.wfPoolFreeCount[1]).To(Equal(16))
		Expect(r.wfPoolFreeCount[2]).To(Equal(6))
		Expect(r.wfPoolFreeCount[3]).To(Equal(6))
		Expect(r.wfPoolFreeCount[4]).To(Equal(6))
	})

	It("should clear reservation when unmap wg", func() {
		wg := kernels.NewWorkGroup()
		for i := 0; i < 10; i++ {
			wf := kernels.NewWavefront(wfSize)
			wg.Wavefronts = append(wg.Wavefronts, wf)
		}
		co.WIVgprCount = 16
		co.WGGroupSegmentByteSize = 1024
		co.WFSgprCount = 64
		wg.CodeObject = co

		r.ReserveResourceForWG(wg)
		r.FreeResourcesForWG(wg)
		assertAllResourcesFree(r)
	})
})
