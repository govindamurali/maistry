package maistry

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan jobChannel
	JobChannel jobChannel
	quit       chan bool
	logger     ILogger
}

func (w *Worker) Start() {

	go func() {
		for {

			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:

				w.logger.Trace("maistry |Worker |Starting Job", nil)
				// we have received a work request.

				job.Do()

				w.logger.Trace("maistry |Worker |Finished Job", nil)

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func newWorker(workerPool chan jobChannel, extLogger ILogger) *Worker {
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(jobChannel),
		quit:       make(chan bool),
		logger:     extLogger,
	}
}
