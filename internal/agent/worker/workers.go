package workerpool

type WorkerPool struct {
	workers chan struct{}
}

func NewWorkerPool(rateLimit int64) *WorkerPool {
	return &WorkerPool{
		workers: make(chan struct{}, rateLimit),
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.workers <- struct{}{}
	go func() {
		task()
		<-wp.workers
	}()
}