package main

import (
  "fmt"
  "./fsm"
)

const (
  GARAGE_STATE_OPENED         = iota
  GARAGE_STATE_CLOSED         = iota
  GARAGE_STATE_OPENING        = iota
  GARAGE_STATE_CLOSING        = iota
  GARAGE_STATE_HALTED_OPENING = iota
  GARAGE_STATE_HALTED_CLOSING = iota
)

const (
  GARAGE_EVENT_BUTTON_PRESSED = iota
  GARAGE_EVENT_BLOCKED        = iota
  GARAGE_EVENT_REACHED_END    = iota
)

func openedCB(id int) {
  fmt.Println("garage opened\n")
}

func closedCB(id int) {
  fmt.Println("garage closed\n")
}

func openingCB(id int) {
  fmt.Println("garage opening\n")
}

func closingCB(id int) {
  fmt.Println("garage closing\n")
}

func haltedCB(id int) {
  fmt.Println("garage halted\n")
}

func haltedOpeningCB(id int) {
  fmt.Println("garage halted while opening\n")
}

func haltedClosingCB(id int) {
  fmt.Println("garage halted while closing\n")
}

func closingTCB(id int, arg interface{}) {
  fmt.Println("garage transition closing\n")
}

func haltedOpeningTCB(id int, arg interface{}) {
  fmt.Println("garage transition halted-opening\n")
}

func haltedClosingTCB(id int, arg interface{}) {
  fmt.Println("garage transition halted-closing\n")
}

func openingTCB(id int, arg interface{}) {
  fmt.Println("garage transition opening\n")
}

func closedTCB(id int, arg interface{}) {
  fmt.Println("garage transition closed\n")
}

func openedTCB(id int, arg interface{}) {
  fmt.Println("garage transition opened\n")
}

var g_states = []fsm.FSMStateType {
  fsm.FSMStateType {
    StateId: GARAGE_STATE_OPENED, 
    StateName: "opened", 
    Callback: openedCB,
  },
  fsm.FSMStateType {
    StateId: GARAGE_STATE_CLOSED, 
    StateName: "closed", 
    Callback: closedCB,
  },
  fsm.FSMStateType {
    StateId: GARAGE_STATE_OPENING, 
    StateName: "opening", 
    Callback: openingCB,
  },
  fsm.FSMStateType {
    StateId: GARAGE_STATE_CLOSING, 
    StateName: "closing", 
    Callback: closingCB,
  },
  fsm.FSMStateType {
    StateId: GARAGE_STATE_HALTED_OPENING, 
    StateName: "halted-opening", 
    Callback: haltedOpeningCB,
  },
  fsm.FSMStateType {
    StateId: GARAGE_STATE_HALTED_CLOSING, 
    StateName: "halted-closing", 
    Callback: haltedClosingCB,
  },
}

var g_events = []fsm.FSMEventType {
  fsm.FSMEventType {
    EventId: GARAGE_EVENT_BUTTON_PRESSED, 
    EventName: "button-pressed", 
  },
  fsm.FSMEventType {
    EventId: GARAGE_EVENT_BLOCKED, 
    EventName: "blocked", 
  },
  fsm.FSMEventType {
    EventId: GARAGE_EVENT_REACHED_END, 
    EventName: "reached-end", 
  },
}

var g_transitions = []fsm.FSMTransitionType {
  fsm.FSMTransitionType {
    State: GARAGE_STATE_CLOSED,
    Event: GARAGE_EVENT_BUTTON_PRESSED,
    Callback: openingTCB,
    NewState: GARAGE_STATE_OPENING,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_OPENING,
    Event: GARAGE_EVENT_BUTTON_PRESSED,
    Callback: haltedOpeningTCB,
    NewState: GARAGE_STATE_HALTED_OPENING,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_HALTED_OPENING,
    Event: GARAGE_EVENT_BUTTON_PRESSED,
    Callback: closingTCB,
    NewState: GARAGE_STATE_CLOSING,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_CLOSING,
    Event: GARAGE_EVENT_BUTTON_PRESSED,
    Callback: haltedClosingTCB,
    NewState: GARAGE_STATE_HALTED_CLOSING,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_HALTED_CLOSING,
    Event: GARAGE_EVENT_BUTTON_PRESSED,
    Callback: openingTCB,
    NewState: GARAGE_STATE_OPENING,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_OPENING,
    Event: GARAGE_EVENT_REACHED_END,
    Callback: openedTCB,
    NewState: GARAGE_STATE_OPENED,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_OPENED,
    Event: GARAGE_EVENT_BUTTON_PRESSED,
    Callback: closingTCB,
    NewState: GARAGE_STATE_CLOSING,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_CLOSING,
    Event: GARAGE_EVENT_REACHED_END,
    Callback: closedTCB,
    NewState: GARAGE_STATE_CLOSED,
  },
  fsm.FSMTransitionType {
    State: GARAGE_STATE_CLOSING,
    Event: GARAGE_EVENT_BLOCKED,
    Callback: openingTCB,
    NewState: GARAGE_STATE_OPENING,
  },
}

func main() {
  fmt.Printf("Hello, world.\n")
  f := fsm.FSMInit("garage", 11, g_states, g_events, g_transitions, GARAGE_STATE_CLOSED, 5)
  f.FSMStart()
  if f.FSMGetState() != GARAGE_STATE_CLOSED {
    fmt.Printf(" Improper State: %s, expedcted %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_CLOSED))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_OPENING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_OPENING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_HALTED_OPENING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_HALTED_OPENING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_CLOSING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_CLOSING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_HALTED_CLOSING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_HALTED_CLOSING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_OPENING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_OPENING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_REACHED_END, nil) 
  if f.FSMGetState() != GARAGE_STATE_OPENED {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_OPENED))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_CLOSING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_CLOSING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_REACHED_END, nil) 
  if f.FSMGetState() != GARAGE_STATE_CLOSED {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_CLOSED))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BUTTON_PRESSED, nil) 
  if f.FSMGetState() != GARAGE_STATE_OPENING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_OPENING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_BLOCKED, nil) 
  if f.FSMGetState() != GARAGE_STATE_OPENING {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_OPENING))
    return
  }
  f.FSMPostEvent(GARAGE_EVENT_REACHED_END, nil) 
  if f.FSMGetState() != GARAGE_STATE_OPENED {
    fmt.Printf(" Improper State: %s, expected %s\n",
      f.FSMGetStateName(f.FSMGetState()), f.FSMGetStateName(GARAGE_STATE_OPENED))
    return
  }
  fmt.Printf("History: \n")
  for _, element := range f.FSMGetHistory() {
    if len(element) > 0 {
      fmt.Printf("%v\n", element)
    }
  }
}
