package moore

import (
	"testing"
	"time"
	"errors"
)

func TestLoop(t *testing.T) {
	success := false
	mm := Make(
		0,
		100,
		func(s State, i Input) (State, error) { return s.(int) + i.(int), nil },
		func() Input { return 1 },
		func(s State) {
			if s.(int) == 100 {
				success = true
			}
		},
	)
	if err := <- mm.Fork(time.NewTicker(time.Duration(time.Nanosecond))); err != nil {
		t.Error("Error in looping")
	}
	if !success {
		t.Error("Expected", true, "got", false)
	}
}

func TestError(t *testing.T) {
	mm := Make(
		0,
		-1,
		func(s State, i Input) (State, error) {
			if s.(int) == 100 {
				return -1, errors.New("I had a failure")
			}
			return s.(int) + i.(int), nil
		},
		func() Input { return 1 },
		func(s State) {},
	)
	if err := <- mm.Fork(time.NewTicker(time.Duration(time.Nanosecond))); err == nil {
		t.Error("Expected error in looping")
	}
}

