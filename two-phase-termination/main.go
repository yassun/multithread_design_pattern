package main

import (
	"fmt"
	"runtime"
	"time"
)

type CountUpThread struct {
	counter           int
	shutdownRequested bool
	ch                chan bool
}

func (ct *CountUpThread) shutdownRequest() {
	ct.shutdownRequested = true
}

func (ct *CountUpThread) join() {
	<-ct.ch
}

func (ct CountUpThread) isShutdownRequested() bool {
	return ct.shutdownRequested
}

func (ct *CountUpThread) run() {

	for !ct.isShutdownRequested() {
		ct.doWork()
	}
	ct.doShutdown()
}

func (ct *CountUpThread) doWork() {
	ct.counter++
	fmt.Printf("doWork: counter = %d\n", ct.counter)
	time.Sleep(500 * time.Millisecond)
}

func (ct *CountUpThread) doShutdown() {
	fmt.Printf("doShutdown: counter = %d\n", ct.counter)
	ct.ch <- true
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("main: BEGIN")

	fmt.Println("main: CountUpThread.run")
	t := CountUpThread{counter: 0, shutdownRequested: false, ch: make(chan bool)}
	go t.run()

	fmt.Println("main: sleep")
	time.Sleep(10000 * time.Millisecond)

	fmt.Println("main: shutdownRequest")
	t.shutdownRequest()

	fmt.Println("main: join")
	t.join()

	fmt.Println("main: END")
}
