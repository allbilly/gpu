package cu

import "gitlab.com/akita/navisim/timing/wavefront"

// An IssueArbiter decides which wavefront can issue instruction
type IssueArbiter struct {
	wfPool *WavefrontPool
}

// NewIssueArbiter returns a newly created IssueArbiter
func NewIssueArbiter(wfPool *WavefrontPool) *IssueArbiter {
	a := new(IssueArbiter)
	a.wfPool = wfPool
	return a
}

// Arbitrate will take a 7 wavefronts in its pool, oldest first.
func (a *IssueArbiter) Arbitrate() []*wavefront.Wavefront {
	if len(a.wfPool.wfs) == 0 {
		return []*wavefront.Wavefront{}
	}

	typeMask := make([]bool, 7)
	list := make([]*wavefront.Wavefront, 0)
	for _, wf := range a.wfPool.wfs {
		if wf.State != wavefront.WfReady || wf.InstToIssue == nil {
			continue
		}

		if !typeMask[wf.InstToIssue.ExeUnit] {
			list = append(list, wf)
			typeMask[wf.InstToIssue.ExeUnit] = true
		}
	}
	return list
}
