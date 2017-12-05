/*
 * Copyright 2017 - Praneet Bachheti
 *
 * Finite State machine (FSM) in golang
 *
 */
package fsm

import (
  "fmt"
  "strconv"
)

type FSMStateCallback func(int)
type FSMTransitionCallback func(int, interface{})

/*
type FSMStateCallback interface {
  StateCallback(id int)
}

type FSMTransitionCallback interface {
  TransitionCallback(id int, arg interface{})
}
*/

type FSMStateType struct {
  StateId    int
  StateName  string
  Callback   FSMStateCallback
}

type FSMEventType struct {
  EventId   int
  EventName string
}

type FSMTransitionType struct {
  State       int
  Event       int
  Callback    FSMTransitionCallback
  NewState    int
}

type FSMType struct {
  name      string
  fsmId     int

  initialized bool

  initState int
  currState int

  stateTable []FSMStateType
  stateCount int

  eventTable []FSMEventType
  eventCount int

  transitionTable []FSMTransitionType
  transitionCount int

  hist_buff []string
  hist_idx  int
  hist_max  int
}

/*-----------------------------------------------------------------------
 * Function : FSMGetHistory
 * Purpose : Get Event/StateChange History
 * Parameters : None
 * Return : array of string
 *-----------------------------------------------------------------------*/
func (f *FSMType) FSMGetHistory() ([]string) {
  return append(f.hist_buff[f.hist_idx:], f.hist_buff[:f.hist_idx]...)
}

/*-----------------------------------------------------------------------
 * Function : FSMAddToHistory
 * Purpose : Add Event/StateChange to History
 * Parameters : currState - old state 
 *              event - event posted
 *              newState - new state to transition to
 * Return : None
 *-----------------------------------------------------------------------*/
func (f *FSMType) FSMAddToHistory(currState int, newState int, event *int) {
  var str string
  if event == nil {
    str = fmt.Sprintf("InitState: %s.", f.FSMGetStateName(newState))
  } else {
    str = fmt.Sprintf("OldState: %s, Event: %s, NewState: %s.",
                       f.FSMGetStateName(currState), f.FSMGetEventName(*event),
                       f.FSMGetStateName(newState))
  }
  f.hist_buff[f.hist_idx] = str
  f.hist_idx++
  if f.hist_idx >= f.hist_max {
    f.hist_idx = 0
  }
}

/*-----------------------------------------------------------------------
 * Function : FSMChangeState
 * Purpose : Moves the state machine to a new state.
 * Parameters : event - pointer to event
 *              newState - new state to transition to
 * Return : 0 is the state machine changed state successfully or an error
 *          code.
 *-----------------------------------------------------------------------*/
func (f *FSMType) FSMChangeState(event *int, newState int) (err int) {
  if newState > f.stateCount {
    f.FSMLog("Illegal state (%d)\n", newState);
    return -1
  }
  if event == nil {
    f.FSMLog("Changing state to %s(%d)\n", f.FSMGetStateName(newState),
             newState);
  } else {
    if newState != f.currState {
      f.FSMLog("Changing state from %s(%d) to %s(%d)\n",
        f.FSMGetStateName(f.currState), f.currState,
        f.FSMGetStateName(newState), newState);
    }
  }

  // Add to FSM history
  f.FSMAddToHistory(f.currState, newState, event)

  /* Change state */
  f.currState = newState

  /* Execute the action associated with the new state */
  cb := f.stateTable[newState].Callback
  if cb != nil {
    f.FSMLog("Executing action for state %s(%d)\n",
      f.FSMGetStateName(newState), newState);
    cb(f.fsmId)
  }
  
  return 0
}

func (f *FSMType) FSMGetStateName(state int) (name string) {
  return f.stateTable[state].StateName
}

func (f *FSMType) FSMGetEventName(event int) (name string) {
  return f.eventTable[event].EventName
}

/*-----------------------------------------------------------------------
 * Function : FSMLog
 * Purpose : Logs the FSM related messages
 * Parameters : message to be logged
 * Return : None
 *-----------------------------------------------------------------------*/
func (f *FSMType) FSMLog(msg string, args ...interface{}) {
  log_msg := "FSM "  + f.name + ":" + strconv.Itoa(f.fsmId) + "#  " + msg
  fmt.Printf(log_msg, args...)
}

/*-----------------------------------------------------------------------
 * Function : FSMGetState
 * Purpose : Returns the current state of the state machine.
 * Parameters : none
 * Return : Current state of the state machine.
 *-----------------------------------------------------------------------*/
func (f *FSMType) FSMGetState() (state int) {
    return f.currState;
}

/*-----------------------------------------------------------------------
 * Function : FSMPostEvent
 * Purpose : Posts an event to a state machine.
 * Parameters : event - An event ID.
 *              arg - void pointer specific to Application
 * Return : 0 is the event was processed successfully or -1 of not.
 *-----------------------------------------------------------------------*/

func (f *FSMType) FSMPostEvent(event int, arg interface{}) (err int) {
  var callback FSMTransitionCallback
  var newState int
  var i int

  f.FSMLog ("Received event %s(%d) in state %s(%d)\n", 
    f.FSMGetEventName(event), event,
    f.FSMGetStateName(f.currState), f.currState);

  /* Loop through and Get the next state from the transition table */
  for i = 0; i < f.transitionCount; i++ {
    if f.transitionTable[i].State == f.currState && 
       f.transitionTable[i].Event == event {
      newState = f.transitionTable[i].NewState
      callback = f.transitionTable[i].Callback
      break
    }
  }

  if i == f.transitionCount {
    /* event not found in the transition table */
    f.FSMLog("Couldn't find event %s(%d) in state %s(%d) in transition table\n",
      f.FSMGetEventName(event), event, f.FSMGetStateName(f.currState),
      f.currState);

    return -1
  }

  /* Execute the action assosiated with the event */
  if callback != nil {
    f.FSMLog("Executing action for event %s(%d)\n",
      f.FSMGetEventName(event), event);
    callback(f.fsmId, arg)
  }

  f.FSMLog("Changing state to %s(%d)\n",
    f.FSMGetStateName(newState), newState);

  return f.FSMChangeState(&event, newState)
}

/*-----------------------------------------------------------------------
 * Function : FSMStart
 * Purpose : Starts executing the Finite State Machine
 * Parameters : none
 * Return : 0 if the Finite State Machine started successfully and -1
 * in case of an error.
 *-----------------------------------------------------------------------*/
func (f *FSMType) FSMStart() (err int) {
  if f.initialized {
    f.FSMLog("FSM Already initialized in State %s\n", f.FSMGetStateName(f.initState))
    return -1
  }
  f.initialized = true

  f.FSMChangeState(nil, f.initState)

  f.FSMLog("Started in State %s\n", f.FSMGetStateName(f.initState))

  return 0
}

/*-----------------------------------------------------------------------
 * Function : FSMInit
 * Purpose : Initialize the Finite State Machine
 * Parameters : 
 *            name
 *            fsmId
 *            states
 *            events
 * Return : FSM Object pointer
 *-----------------------------------------------------------------------*/
func FSMInit(name string, fsm_id int, states []FSMStateType,
             events []FSMEventType, transitions []FSMTransitionType,
             init_state int, max_history int) (fsm *FSMType) {
  fsm = new(FSMType)

  fsm.name = name
  fsm.fsmId = fsm_id

  fsm.initState = init_state
  fsm.currState = init_state

  fsm.stateTable = states
  fsm.stateCount = len(states)

  fsm.eventTable = events
  fsm.eventCount = len(events)

  fsm.transitionTable = transitions
  fsm.transitionCount = len(transitions)

  fsm.hist_buff = make([]string, max_history)
  fsm.hist_idx = 0
  fsm.hist_max = max_history

  fsm.FSMLog("initialized state Count (%d), event Count (%d)",
             " Trans Count (%d)\n", fsm.stateCount, fsm.eventCount,
             fsm.transitionCount)

  return fsm
}

