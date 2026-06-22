package cu

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/navisim/timing/wavefront"
)

var _ = Describe("FetchArbiter", func() {
	var (
		wfPool  *WavefrontPool
		arbiter *FetchArbiter
	)

	BeforeEach(func() {

		wfPool = NewWavefrontPool(10)
		arbiter = NewFetchArbiter(wfPool)
		arbiter.InstBufByteSize = 256
	})

	It("should find the oldest wf to dispatch", func() {
		wfLastFetchTime := []akita.VTimeInSec{
			10.2, 10.3, 9.8, 9.7, 9.4,
			9.6, 9.5, 9.6, 9.8, 10.0,
		}
		wfState := []wavefront.WfState{
			wavefront.WfRunning, wavefront.WfRunning, wavefront.WfReady, wavefront.WfReady, wavefront.WfRunning,
			wavefront.WfRunning, wavefront.WfRunning, wavefront.WfRunning, wavefront.WfReady, wavefront.WfRunning,
		}

		for i := 0; i < len(wfState); i++ {
			wf := new(wavefront.Wavefront)
			wf.LastFetchTime = wfLastFetchTime[i]
			wf.State = wfState[i]
			wfPool.AddWf(wf)

			if i == 4 {
				wf.InstBuffer = make([]byte, arbiter.InstBufByteSize)
			}
		}

		wfs := arbiter.Arbitrate()

		Expect(len(wfs)).To(Equal(1))
		Expect(wfs[0].LastFetchTime).To(Equal(akita.VTimeInSec(9.5)))
	})
})
