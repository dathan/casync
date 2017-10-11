package casync

import "time"


/**
	A task executes
	We need to tell tasks if there is an error shutdown and don't execute
	We need to interrupt a task which is running if a signal tells us to or the Tasks is executing to long
	We need to control the number of executing tasks
	We want to enter a task for executing when the task is ready to execute
	We need to know when all tasks have completed.
 */
//
// Task to execute, each task has a timeout in seconds
//
type Task struct {
	id int
	Execute func() // what to execute
	// todo add how long A task should last
	timeout_sec time.Duration
}

//
// Build the Task
//
func NewTask(id int, exec func()) *Task {


	return &Task {
		id : id,
		Execute: exec,
		timeout_sec: 0,
	}
}

//
// set a timeout
//
func(t *Task) SetTimeout(s int) {
	t.timeout_sec = time.Duration(s) * time.Second
}