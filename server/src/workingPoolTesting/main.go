package main

import (
	"./workers"
	"fmt"
)


type mytask struct {
	str string
}

func (t mytask) Handle() error {
	fmt.Println("Done", t.str)
	return nil
}

func (t mytask) Finish(err error) {
	fmt.Println(err)
}

func main() {
	wp := workers.NewPool(10)


	defer wp.Stop()

	tsk := mytask{"It's like a piece of cake."}

	wp.AddTaskAsynk(tsk,1)

	wp.Run()

}
