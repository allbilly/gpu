package cu

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/navisim/timing/wavefront"
)

var _ = Describe("IssueArbiter", func() {
	var (
		arbiter *IssueArbiter
		wfPool  *WavefrontPool
	)

	BeforeEach(func() {

		wfPool = NewWavefrontPool(10)
		arbiter = NewIssueArbiter(wfPool)

	})

	It("should decide which wf to issue", func() {
		wfState := []wavefront.WfState{
			wavefront.WfRunning, wavefront.WfReady, wavefront.WfReady, wavefront.WfReady, wavefront.WfReady,
			wavefront.WfReady, wavefront.WfReady, wavefront.WfReady, wavefront.WfReady, wavefront.WfReady,
		}
		exeUnits := []rdnainsts.ExeUnit{
			rdnainsts.ExeUnitVALU, rdnainsts.ExeUnitScalar, rdnainsts.ExeUnitVMem,
			rdnainsts.ExeUnitBranch, rdnainsts.ExeUnitLDS, rdnainsts.ExeUnitSpecial,
			rdnainsts.ExeUnitVALU, rdnainsts.ExeUnitBranch, rdnainsts.ExeUnitVALU,
			rdnainsts.ExeUnitVMem,
		}
		wfs := make([]*wavefront.Wavefront, 0)

		for i := 0; i < len(wfState); i++ {
			wf := new(wavefront.Wavefront)
			wf.State = wfState[i]
			wf.InstToIssue = wavefront.NewInst(rdnainsts.NewInst())
			wf.InstToIssue.ExeUnit = exeUnits[i]
			wfs = append(wfs, wf)
			wfPool.AddWf(wf)

			if i == 3 || i == 6 {
				wf.InstToIssue = nil
			}
		}

		issueCandidate := arbiter.Arbitrate()

		Expect(len(issueCandidate)).To(Equal(6))
		Expect(issueCandidate).NotTo(ContainElement(BeIdenticalTo(wfs[0])))
		Expect(issueCandidate).To(ContainElement(BeIdenticalTo(wfs[1])))
		Expect(issueCandidate).To(ContainElement(BeIdenticalTo(wfs[2])))
		Expect(issueCandidate).NotTo(ContainElement(BeIdenticalTo(wfs[3])))
		Expect(issueCandidate).To(ContainElement(BeIdenticalTo(wfs[4])))
		Expect(issueCandidate).To(ContainElement(BeIdenticalTo(wfs[5])))
		Expect(issueCandidate).NotTo(ContainElement(BeIdenticalTo(wfs[6])))
		Expect(issueCandidate).To(ContainElement(BeIdenticalTo(wfs[7])))
		Expect(issueCandidate).To(ContainElement(BeIdenticalTo(wfs[8])))
		Expect(issueCandidate).NotTo(ContainElement(BeIdenticalTo(wfs[9])))
	})
})
