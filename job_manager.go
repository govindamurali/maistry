package maistry

import (
	"time"
)

type JobManager struct {
	// A pool of workers channels that are registered with the dispatcher
	allocatedWorkerChan chan jobChannel

	// Used to throttle
	throttler <-chan time.Time

	// To communicate from PushJob() To Run()
	jobChan jobChannel

	MaxJps int
	name   string
}

func (jm *JobManager) Run() {
	go func() {
		for {
			// Throttle the loop according to input Jps
			<-jm.throttler

			// Wait for a job from PushJob()
			job := <-jm.jobChan

			// Get a worker (type JobChannel) from Worker Pool
			worker := <-jm.allocatedWorkerChan

			// Push the job to job channel
			worker <- job

		}
	}()
}

// PushJob - blocking call to JobManager
func (jm *JobManager) PushJob(j Job) {
	jm.jobChan <- j
}

// NewJobManager - jps (Jobs per second) - JobManager will throttle job execution according if it crosses maxJps
func NewJobManager(maxJps int, name string) *JobManager {

	throttler := time.Tick(time.Duration(int(1000/maxJps)) * time.Millisecond)

	return &JobManager{
		// Zero size WorkerPool
		allocatedWorkerChan: make(chan jobChannel),
		// To throttle
		throttler: throttler,
		// For communicating from PushJob() To Run()
		jobChan: make(jobChannel),
		MaxJps:  maxJps,
		name:    name,
	}
}
