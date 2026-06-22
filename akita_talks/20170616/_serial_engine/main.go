package main

import (
	"log"
	"time"

	"gitlab.com/yaotsu/core"
)

// START OMIT
func main() {
	engine := core.NewSerialEngine()
	handler := new(exampleHandler)

	for i := 0; i < 1000; i++ {
		evt := newTestEvent(1.0, handler)
		engine.Schedule(evt) // HL
	}

	start := time.Now()
	engine.Run() // HL
	elapsed := time.Since(start)
	log.Printf("Simulation takes %s", elapsed)
}

// END OMIT

type testEvent struct {
	*core.EventBase
}

func newTestEvent(t core.VTimeInSec, handler core.Handler) *testEvent {
	evt := new(testEvent)
	evt.EventBase = core.NewEventBase(t, handler)
	return evt
}

type exampleHandler struct {
}

func (h *exampleHandler) Handle(evt core.Event) error {
	for i := 0; i < 10000000; i++ {
	}
	// log.Print("Complete\n")
	return nil
}
