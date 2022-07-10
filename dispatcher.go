package maistry

import (
	"reflect"
	"sync"
)

type Dispatcher struct {
	jobManagers []*JobManager
	workerPool  chan jobChannel
	maxWorkers  int
	sync.RWMutex

	logger                 ILogger
	jobManagerCommunicator chan *JobManager
}

func NewDispatcher(maxWorkers int, extLogger ILogger) *Dispatcher {

	return &Dispatcher{
		maxWorkers:             maxWorkers,
		workerPool:             make(chan jobChannel, maxWorkers),
		jobManagerCommunicator: make(chan *JobManager),
		logger:                 extLogger,
	}
}

//This Will Add a job manager to the list
func (d *Dispatcher) AddJobManager(jm *JobManager) {
	d.Lock()
	defer d.Unlock()
	d.jobManagers = append(d.jobManagers, jm)

	go func() {
		d.jobManagerCommunicator <- jm
	}()
}

func (d *Dispatcher) Start() (started bool) {

	d.startCheck()

	d.startWorkers(d.maxWorkers)

	go d.dispatch()

	return true
}

func (d *Dispatcher) dispatch() {

	var selectCases []reflect.SelectCase

	// equivalent to - case jobManager := <-d.jobManagerCommunicator
	newJobManagerReceivedCase := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(d.jobManagerCommunicator)}
	selectCases = append(selectCases, newJobManagerReceivedCase)

	for {
		// A worker got free
		worker := <-d.workerPool

		// Starting from i = 1, skipping the 1st selectCase as its the newJobManagerReceivedCase
		// we don't send incoming worker to newJobManagerReceivedCase
		for i := 1; i < len(selectCases); i++ {

			// equivalent to - case jobManager.allocatedWorkerChan <- worker
			selectCases[i].Send = reflect.ValueOf(worker)
		}

		chosen, recvVal, recvOK := reflect.Select(selectCases)

		// Check if newJobManagerReceivedCase was successful
		if chosen == 0 && recvOK {

			// send worker back to the pool
			d.workerPool <- worker

			// Get newly received JobManager from reflect.Value
			jobManager, ok := recvVal.Interface().(*JobManager)
			if !ok {
				panic("recevied value is not jobManager")
			}

			// Add the newly received jobManager's allocatedWorkerChan to the selectCase
			// equivalent to - case jobManager.allocatedWorkerChan <- XXXXX
			jobManagerAllocateWorkerCase := reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(jobManager.allocatedWorkerChan)}
			selectCases = append(selectCases, jobManagerAllocateWorkerCase)
		}
	}
}

func runJobManagers(jobManagers []*JobManager) {

	for _, jm := range jobManagers {
		jm.Run()
	}
}

func (d *Dispatcher) startWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		w := newWorker(d.workerPool, d.logger)
		w.Start()
	}
}

func (d *Dispatcher) startCheck() {
	if d.maxWorkers <= 0 {
		panic("Max worker configured <= 0")
	}

	if len(d.jobManagers) == 0 {
		panic("No JobManagers added")
	}
}
