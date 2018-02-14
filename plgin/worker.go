/*
 * Copyright 2018 - Praneet Bachheti
 *
 * Dispatcher Implementation
 *
 */

package main

import (
  "fmt"
  "time"
)

type Gofer struct {
  GoferID      int
  JobChan      chan JobRequest
  ExitChan     chan bool
  GoferQueue   chan chan JobRequest
}

// Create a New Gofer
// - Create a JobRequest Channel to listen on
// - Create an ExitChan to termnate
// - Add self to the GoferQueue so we get JobRequests
func CreateGofer(id int, goferQueue chan chan JobRequest) Gofer {

  gofer := Gofer {
    GoferID: id,
    JobChan: make(chan JobRequest),
    GoferQueue: goferQueue,
    ExitChan: make(chan bool),
  }

  return gofer
}

func (g *Gofer) Run() {
  go func() {
    for {
      // Be part of the Gofer Queue
      g.GoferQueue <- g.JobChan

      select {

      case job := <-g.JobChan:
        // Received a Job Request, process it
        fmt.Printf("Gofer: %d, Received job request %d\n", g.GoferID, job.JobID)
        time.Sleep(5 * time.Second)
        fmt.Printf("Gofer: %d, Slept for 5 seconds\n", g.GoferID)

      case <-g.ExitChan:
        fmt.Printf("Gofer: %d exiting\n", g.GoferID)
        return

      }
    }
  }()
}


func (g *Gofer) Exit() {
  go func() {
    g.ExitChan <- true
  }()
}
