package moore

import (
	"errors"
	"reflect"
	"runtime"
	"time"
)

type State interface{}
type Input interface{}

type TransitionFunction func(State, Input) (State, error)
type OutputFunction func(State)
type InputFunction func() Input

type MooreMachine struct {
	currentState       State
	quitState          State
	transitionFunction TransitionFunction
	inputFunction      InputFunction
	outputFunction     OutputFunction
}

func Make(startState State, quitState State, transitionFunction TransitionFunction, inputFunction InputFunction, outputFunction OutputFunction) MooreMachine {
	return MooreMachine{
		currentState:       startState,
		quitState:          quitState,
		transitionFunction: transitionFunction,
		inputFunction:      inputFunction,
		outputFunction:     outputFunction,
	}
}

func (mm *MooreMachine) Fork(ticker *time.Ticker) chan error {
	errorChannel := make(chan error) // Use a channel to pass any error back to user and allow them to wait until quit state is reached.
	go func() {
		runtime.LockOSThread()
		errorChannel <- mm.Run(ticker)
	}()
	return errorChannel
}

func (mm *MooreMachine) Run(ticker *time.Ticker) error {
	for range ticker.C { // Loop based on a timer.
		var err error
		if err = mm.Verify(); err != nil { // Verify that variable types are correct.
			return err
		}
		if mm.currentState == mm.quitState || reflect.DeepEqual(mm.currentState, mm.quitState) { // Quit if this is the quit state.
			return nil // No error encountered.
		}

		mm.currentState, err = mm.transitionFunction(mm.currentState, mm.inputFunction()) // Do a state transition.

		if mm.currentState != nil {
			mm.outputFunction(mm.currentState) // Do output for current state.
		}
		if err != nil {                    // Send error and quit.
			return err
		}
	}
	return nil
}

func (mm *MooreMachine) Verify() error {
	stateType := reflect.TypeOf(mm.currentState)
	if stateType != nil {
		if mm.quitState != nil && !stateType.AssignableTo(reflect.TypeOf(mm.quitState)) {
			return errors.New("type of current state differs from that of quit state")
		}
		if !stateType.AssignableTo(reflect.TypeOf(mm.transitionFunction).In(0)) {
			return errors.New("type of current state differs from that of transition function argument")
		}
		if !stateType.AssignableTo(reflect.TypeOf(mm.transitionFunction).Out(0)) {
			return errors.New("type of current state differs from that of transition function return")
		}
		if !stateType.AssignableTo(reflect.TypeOf(mm.outputFunction).In(0)) {
			return errors.New("type of current state differs from that of output function argument")
		}
	}

	inputType := reflect.TypeOf(mm.inputFunction).Out(0)
	if !inputType.AssignableTo(reflect.TypeOf(mm.transitionFunction).In(1)) {
		return errors.New("type of input function return differs from that of transition function argument")
	}
	return nil
}
