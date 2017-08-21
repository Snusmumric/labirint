package workers

import (
	"time"
)

type Task interface {
	Handle() error
	Finish(error)
}


type WorkerPool struct {
	capacity int
	taskChan chan Task
}

func NewPool(capacity int) *WorkerPool{
	return &WorkerPool{ capacity: capacity, taskChan: make(chan Task, capacity)}
}

func (p *WorkerPool) Run() {
	for i := 0; i < p.capacity; i++ {
		go p.startWorker()
	}
}

func (p *WorkerPool) Stop() {
	close(p.taskChan)
}

func (p *WorkerPool) AddTaskAsynk(task Task, timeout time.Duration) error {
	//tick := time.Tick(timeout * time.Millisecond)
	p.taskChan <- task
	return nil
	/*
	select {
	case p.taskChan <- task:
		return nil
	case <-tick:
		return fmt.Errorf("Task not sended, timeout")
	}
	*/


}

func (p *WorkerPool) startWorker() {
	for t := range p.taskChan {
		err := t.Handle()
		t.Finish(err)
	}
}
