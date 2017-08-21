package worker_pool

import (
	"fmt"
	"time"
)

type Task interface {
	Handle() error
	Finish(error)
}

type WorkerPool struct {
	capacity int
	taskChan chan *Task
}

func NewPool(capacity int) {
	return &WorkerPool{
		capacity: capacity,
		taskChan: make(chan *Task),
	}
}

func (p *WorkerPool) Run() {
	for i := 0; i < p.capacity; i++ {
		go p.startWorker()
	}
}

func (p *WorkerPool) Stop() {
	close(p.taskChan)
}

func (p *WorkerPool) AddTaskAsynk(task *Task, timeout time.Millisecond) error {
	tick := time.Tick(timeout * time.Millisecond)
	select {
	case p.taskChan <- &Task:
		return nil
	case <-tick:
		return fmt.Errorf("Task not sended, timeout")
	}

}

func (p *Pool) startWorker() {
	for t := range p.taskChan {
		err := t.Handle()
		t.Finish(err)
	}
}
