package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 10 * time.Second

	for i := 0; i < 10; i++ {
		dinnerFinished = []string{}
		dine()
		if len(dinnerFinished) != 5 {
			t.Errorf("incorrect length of slice, Expected 5 but found %d", len(dinnerFinished))
		}
	}
}

func Test_dineWithVaryingDelays(t *testing.T) {
	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", time.Second * 0},
		{"quater secnod delay", time.Microsecond * 250},
		{"half secnod delay", time.Microsecond * 500},
	}

	for _, e := range theTests {
		dinnerFinished = []string{}

		eatTime = e.delay
		sleepTime = e.delay
		thinkTime = e.delay

		dine()
		if len(dinnerFinished) != 5 {
			t.Errorf("%s: incorrect length of slice; expected 5 but got %d", e.name, len(dinnerFinished))
		}
	}
}
