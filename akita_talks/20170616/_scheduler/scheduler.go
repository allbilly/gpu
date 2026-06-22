package scheduler

import (
	"gitlab.com/yaotsu/core"
)

type Scheduler struct {
	*core.ComponentBase

	engine core.Engine
	Freq   core.Freq

	InstMem          core.Component
	SRegFile         core.Component
	BranchUnit       core.Component
	VectorMemDecoder core.Component
	ScalarDecoder    core.Component
	VectorDecoder    core.Component
	LDSDecoder       core.Component
}
