package main

import (
	"fmt"
	"log"
	"time"
)

type Payload struct {
	data string
}

func (p *Payload) UploadData() error {
	time.Sleep(1 * time.Second)
	fmt.Println("work done")
	return nil
}

type Job struct {
	Payload Payload
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				if err := job.Payload.UploadData(); err != nil {
					log.Printf("Error uploading to S3: %s", err.Error())
				}
				log.Println("Success data")
			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
