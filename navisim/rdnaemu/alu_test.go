package rdnaemu

import (
	"gitlab.com/akita/navisim/rdnainsts"
	"gitlab.com/akita/util/ca"
)

type mockInstState struct {
	inst       *rdnainsts.Inst
	scratchpad Scratchpad
}

func (s *mockInstState) PID() ca.PID {
	return 1
}

func (s *mockInstState) Inst() *rdnainsts.Inst {
	return s.inst
}

func (s *mockInstState) Scratchpad() Scratchpad {
	return s.scratchpad
}
