/*
 * Copyright 2018 - Praneet Bachheti
 *
 * Dispatcher Implementation
 *
 */

package main

import (
  "fmt"
)

// gofer: a person whose job is to do various small jobs.

// Each Gofer has an unbuffered channel to listen on
// GoferQueue is a buffered channel of these channel type
var GoferQueue chan chan JobRequest

func InitDispatcher(numGofers int) {

  // Initialize the GoferQueue
  GoferQueue = make(chan chan JobRequest, numGofers)

  // Create the Gofers and Run them
  for idx := 0; idx < numGofers; idx++ {

    gofer := CreateGofer(idx+1, GoferQueue)

    gofer.Run()

    fmt.Println("Started Gofer", idx+1)

  }

}

func AssignJob(job JobRequest) {

  // Get an Idle Gofer
  gofer := <-GoferQueue

  // Assign Gofer the Job to work on
  gofer <- job

  fmt.Printf("Assigned Job: %d to Gofer\n", job.JobID)

}

func RunDispatcher() {

  fmt.Println("Run Dispacther")

  go func() {
    // Run Forever
    for {
      select {
      case job := <-JobQueue:
        // Assign the Job to a Gofer
        go AssignJob(job)
      }
    }
  }()

}
