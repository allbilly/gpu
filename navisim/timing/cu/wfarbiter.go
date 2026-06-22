package cu

import "gitlab.com/akita/navisim/timing/wavefront"

// An WfArbiter can decide which wavefront can take action,
// in a list of wavefront pools
type WfArbiter interface {
	Arbitrate() []*wavefront.Wavefront
}
