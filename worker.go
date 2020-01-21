package grpool

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	jobChan   <-chan interface{}
	closeChan chan struct{}
	job       Job
}

func newWorker(jobChan chan interface{}, newJob func() Job) Worker {
	w := &worker{
		jobChan:   jobChan,
		job:       newJob(),
		closeChan: make(chan struct{}),
	}
	return w
}

func (w *worker) Start() {
	for {
		select {
		case value := <-w.jobChan:
			w.job.Process(value)
		case <-w.closeChan:
			return
		}
	}
}

func (w *worker) Stop() {
	close(w.closeChan)
}
