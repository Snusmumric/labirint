package worker_pool

type taskArgs []interface{}

type handler func(...interface{})

type Task struct {
	taskArgs taskArgs
	handler  interface{}
}

type WorkerPool struct {
	capacity int
	taskChan chan *Task
}

func NewPool(capacity int, chanSize int) {
	return &WorkerPool{
		capacity: capacity,
		taskChan: make(chan *Task, chanSize),
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

func (p *WorkerPool) AddTaskAsynk(task *Task) {
	p.taskChan <- &Task
}

func (p *Pool) startWorker() {
	for t := range p.taskChan {
		t.handler(t.taskArgs)
	}
}
