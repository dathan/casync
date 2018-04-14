package casync

import (
	"fmt"
	"math/rand"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSignals(t *testing.T) {

	var Tasks []*Task

	for i := 0; i < 10000; i++ {
		t := fakeTask(i)
		Tasks = append(Tasks, NewTask(i, t))
	}

	as := NewAsync(4, Tasks)

	go func() { // after 3 second send a sigterm

		time.Sleep(time.Second * 3)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()

	as.ExecuteTasks()
	//panic("Show stack")

}

func TestSignalsTimeouts(t *testing.T) {

	var Tasks []*Task

	for i := 0; i < 10000; i++ {
		t := fakeTask(i)
		tsk := NewTask(i, t)
		tsk.SetTimeout(2)
		Tasks = append(Tasks, tsk)
	}

	as := NewAsync(4, Tasks)

	go func() { // after 3 second send a sigterm

		time.Sleep(time.Second * 10)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(syscall.SIGTERM)
	}()

	as.ExecuteTasks()
	//panic("Show stack")

}

func TestTimeout(t *testing.T) {

	var Tasks []*Task

	for i := 0; i < 10000; i++ {
		t := fakeTask(i)
		tsk := NewTask(i, t)
		tsk.SetTimeout(1)
		Tasks = append(Tasks, tsk)
	}

	as := NewAsync(4, Tasks)

	as.ExecuteTasks()
	//panic("Show stack")

}

func fakeTask(i int) func() {
	return func() {
		fmt.Printf("Start - Job: %d\n", i)
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		fmt.Printf("DONE - Job: %d\n", i)

	}

}
