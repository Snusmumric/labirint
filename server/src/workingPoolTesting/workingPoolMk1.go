package workPool

import (
	"time"
	"fmt"
)

type Task interface {
	Handle() error
	//Finish(error)
}

type WorkerPool struct {
	capacity int
	taskChan chan Task
}

func NewPool(capacity int) *WorkerPool {
	return &WorkerPool{capacity: capacity, taskChan: make(chan Task)}
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

	//p.taskChan <- task ; return nil


	tick := time.Tick(timeout * time.Microsecond)
	select {
	case p.taskChan <- task:
		return nil
	case <-tick:
		// в случае, если воркеров не хватает, они заняты,
		// то через некоторое время таск будет отклонен.
		// типо попробуйте позднее
		fmt.Println("Task not sended, timeout")
		return fmt.Errorf("Task not sended, timeout")
	}

}

func (p *WorkerPool) startWorker() {
	for t := range p.taskChan {
		_ = t.Handle()
		//t.Finish(err)
	}
}
