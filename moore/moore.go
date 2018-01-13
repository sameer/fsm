package moore

import (
	"time"
	"reflect"
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
	quit               bool
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
	mm.Verify()
	errchan := make(chan error)
	go func() {
		for range ticker.C {
			mm.outputFunction(mm.currentState)
			if mm.currentState == mm.quitState {
				break
			}
			var err error
			mm.currentState, err = mm.transitionFunction(mm.currentState, mm.inputFunction())
			if err != nil {
				errchan <- err
				return
			}
		}
		errchan <- nil
	}()
	return errchan
}

func (mm *MooreMachine) Verify() {
	stateType := reflect.TypeOf(mm.currentState)
	if !stateType.AssignableTo(reflect.TypeOf(mm.quitState)) {
		panic("Type of current state differs from that of quit state.")
	}
	if !stateType.AssignableTo(reflect.TypeOf(mm.transitionFunction).In(0)) {
		panic("Type of current state differs from that of transition function argument.")
	}
	if !stateType.AssignableTo(reflect.TypeOf(mm.transitionFunction).Out(0)) {
		panic("Type of current state differs from that of transition function return.")
	}
	if !stateType.AssignableTo(reflect.TypeOf(mm.outputFunction).In(0)) {
		panic("Type of current state differs from that of output function argument.")
	}

	inputType := reflect.TypeOf(mm.inputFunction).Out(0)
	if !inputType.AssignableTo(reflect.TypeOf(mm.transitionFunction).In(1)) {
		panic("Type of input function return differs from that of transition function argument.")
	}
}
