package main

import (
	"./workPool"
	"fmt"
	"time"
)

type myTask struct {
	str string
}

func (t myTask) Handle() error {
	time.Sleep(1*time.Second)
	fmt.Println(t.str)
	done<- struct{}{}
	return nil
}

func (t myTask) Finish(err error)  {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

	var done = make(chan struct{})

func main() {

	wp := workPool.NewPool(2) // number of workers rules as threads
	wp.Run()

	defer wp.Stop()

	t1 := myTask{"Hello"}
	t2 := myTask{", my dear"}
	t3 := myTask{" fellow."}

	taskList := []myTask{t1,t2,t3}

	for _, t := range taskList {
		go func(t myTask) {
			_ = wp.AddTaskAsynk(t, 1)
			/*
			if err != nil {
				fmt.Println(err)
				wp.Stop()
				return
			}
			*/
		}(t)
	}

	for i:=0; i<3; i++{
		<-done
	}
}
