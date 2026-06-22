package cu

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/navisim/timing/wavefront"
)

// A SubComponent is an element installed in the compute unit
type SubComponent interface {
	CanAcceptWave() bool
	AcceptWave(wave *wavefront.Wavefront, now akita.VTimeInSec)
	Run(now akita.VTimeInSec) bool
	IsIdle() bool
	Flush()
}
