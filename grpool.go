package grpool

type Pool interface {
	Process(interface{})
	Close()
}

type pool struct {
	workers []Worker
	jobChan chan interface{}
}

func New(n int, newJob func() Job) Pool {
	p := &pool{
		jobChan: make(chan interface{}),
	}
	for i := 0; i < n; i++ {
		w := newWorker(p.jobChan, newJob)
		go w.Start()
		p.workers = append(p.workers, w)
	}
	return p
}

func (gp *pool) Process(value interface{}) {
	gp.jobChan <- value
}

func (gp *pool) Close() {
	for _, w := range gp.workers {
		w.Stop()
	}
	close(gp.jobChan)
}