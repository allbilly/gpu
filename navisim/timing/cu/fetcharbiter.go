package cu

import (
	"math"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/navisim/timing/wavefront"
)

// A FetchArbiter can decide which wavefront in a scheduler can fetch
// instructions
type FetchArbiter struct {
	wfPool          *WavefrontPool
	InstBufByteSize int
}

// NewFetchArbiter creates a fetch arbiter
func NewFetchArbiter(wfPool *WavefrontPool) *FetchArbiter {
	a := new(FetchArbiter)
	a.wfPool = wfPool
	return a
}

// Arbitrate decide which wavefront can fetch the next instruction
func (a *FetchArbiter) Arbitrate() []*wavefront.Wavefront {
	list := make([]*wavefront.Wavefront, 0, 1)

	oldestTime := akita.VTimeInSec(math.MaxFloat64)
	var toFetch *wavefront.Wavefront

	for _, wf := range a.wfPool.wfs {
		wf.RLock()

		if wf.IsFetching {
			wf.RUnlock()
			continue
		}

		if len(wf.InstBuffer) >= a.InstBufByteSize {
			wf.RUnlock()
			continue
		}

		if wf.LastFetchTime < oldestTime {
			toFetch = wf
			oldestTime = wf.LastFetchTime
		}
		wf.RUnlock()
	}

	if toFetch != nil {
		list = append(list, toFetch)
	}

	return list
}
