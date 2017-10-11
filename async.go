package casync

import (
	"sync"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
)

type Async struct {
	TasksChannel chan *Task // tasks to execute
	Concurrency  int // how many go routines should be active.
	wg           sync.WaitGroup // control the go routines
	Tasks        []*Task // list of tasks for this async struct
	Sigs         chan os.Signal // catch unix sigs to graceful shutdown
	stop         bool // tell the app to stop
}

//
// build up the async process
//
func NewAsync(sizeofCurrencency int, ts []*Task) *Async {

	as := &Async{
		TasksChannel: make(chan *Task, len(ts)), // make a channel to a pointer of task
		Concurrency: sizeofCurrencency,
		Tasks: ts,
		Sigs : make(chan os.Signal, 1),
		stop : false,
	}

	signal.Notify(as.Sigs, syscall.SIGINT, syscall.SIGTERM) // let as.sigs know if sigtem or sigint happened

	as.setupWorkers()
	go as.WatchTasks() // todo put this in biz logic pull out of the creation

	return as

}


//
// allways run size of concurrency Tasks. For instance if Async.Concurrency == 4, 4 gouroutines will dequeue Tasks to
// execute
//
func (a *Async) setupWorkers() {

	// allow N goroutines to execute at once
	for j := 0; j < a.Concurrency; j++ {

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			for task := range a.TasksChannel {
				// read off all the tasks in each routine

				if a.stop == false {
					// do not execute any more tasks
					if task.timeout_sec > 0 {

						ch := make(chan bool, 1)
						go func() {
							task.Execute()
							ch <- true
						}()

						select {
						case <-ch:
							fmt.Println("Finish: ", task.id)
							continue
						case <-time.After(task.timeout_sec):
							fmt.Println("TIMEOUT: ", task.id)
							continue
						}

					} else {
						task.Execute()
					}

				} else {
					return
				}
			}
		}()
	}
}

//
// now add the tasks to the TasksChannel
//
func (a *Async) ExecuteTasks() {

	for _, item := range a.Tasks {
		a.TasksChannel <- item
	}

	close(a.TasksChannel)
	a.wg.Wait() // block until the N go routines stop

}


//
// catch a signal and drain the tasks
//
func (a *Async) WatchTasks() {

	for {
		select {
		case <-a.Sigs:
		// don't care what the signal is we will just stop the async process
			a.stop = true;
			return // break out of the loop and Exit

		}
	}
}