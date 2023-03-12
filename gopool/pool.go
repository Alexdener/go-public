package pool

import (
	"runtime"
	"sync"
)

type Pool struct {
	WorkerChan   chan *Task
	MaxWorkerNum int
	queue        []*Task
	sync.Mutex
	wg sync.WaitGroup
}

type Task struct {
	f func() error
}

func NewPool(worker int) *Pool {
	pool := &Pool{
		WorkerChan:   make(chan *Task, 2*worker),
		queue:        make([]*Task, worker),
		MaxWorkerNum: worker,
	}
	go pool.run()
	return pool
}

func (p *Pool) Wait() {
	p.waitTask()
	close(p.WorkerChan)
	p.wg.Wait()
	select {
	default:
		return
	}
}

func (p *Pool) waitTask() {
	for {
		runtime.Gosched()
		if len(p.WorkerChan) == 0 && len(p.queue) == 0 {
			break
		}
	}
}

func (p *Pool) run() {
	go func() {
		for {
			if len(p.queue) > 0 {
				p.Lock()
				w := p.queue[0]
				p.queue = p.queue[1:]
				p.WorkerChan <- w
				p.Unlock()
			}
		}
	}()
	p.wg.Add(p.MaxWorkerNum)
	for i := 0; i < p.MaxWorkerNum; i++ {
		go p.worker(i)
	}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for w := range p.WorkerChan {
		if w == nil {
			continue
		}
		w.work()
	}
}

func (p *Pool) NewTask(f func() error) {
	t := &Task{
		f: f,
	}
	p.Lock()
	p.queue = append(p.queue, t)
	p.Unlock()
}

func (t *Task) work() {
	_ = t.f()
}
