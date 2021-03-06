package main

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	NumberOfWorkers int
	Task Task
}

type Worker struct {
	Wg *sync.WaitGroup
	Task Task
}

type Task struct {
	Input chan int // Update `Input` to be required channel type.
	Execute func(n int) // Update `Execute` to be required "Work" function signature.
}

func (wp *WorkerPool) Run (done chan bool) {
	var wg sync.WaitGroup
	wg.Add(wp.NumberOfWorkers)

	for i := 0; i < wp.NumberOfWorkers; i++ {
		w := Worker{
			Wg: &wg,
			Task: wp.Task,
		}

		go w.Work()
	}

	wg.Wait()
	done <- true
}

func (w *Worker) Work() {
	defer w.Wg.Done()

	for i := range w.Task.Input {
		w.Task.Execute(i)
	}
}

func main() {
	input := make(chan int) // Update channel to required channel type.

	done := make(chan bool)
	wp := WorkerPool{
		NumberOfWorkers: 50, // Update number of workers desired.
		Task: Task{
			Input: input,
			Execute: workFunction, // Update target work function.
		},
	}
	go wp.Run(done)

	// Fetch the data needed for processing and send to the Input channel of the WorkPool.
	fetchDataForProcessing(input)

	<-done
}

func fetchDataForProcessing(input chan int) {
	for rs := 0; rs < 200; rs++ {
		input <- rs
	}

	close(input)
}

// workFunction signature should match the individual element in the `input` channel type.
func workFunction(o int) {
	// Simulate work being done.
	time.Sleep(time.Second)
	fmt.Println(o)
}
