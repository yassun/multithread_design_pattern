package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Request struct {
	name string
}

func (r Request) String() string {
	return fmt.Sprintf("[ Request %s ]", r.name)
}

type RequestQueue struct {
	queue *list.List
}

func (rq RequestQueue) getRequest() Request {
	if rq.queue.Len() <= 0 {
		// TODO wait
	}
	e := rq.queue.Back()
	request := e.Value.(Request)
	return request
}

func (rq RequestQueue) putRequest(request Request) {
	rq.queue.PushFront(request)
	// TODO notifyAll
}

type ClientThread struct {
	random       *rand.Rand
	requestQueue RequestQueue
}

func (ct ClientThread) run() {
	for i := 0; i < 10000; i++ {
		request := Request{name: fmt.Sprintf("No. %d", i)}
		fmt.Println("ClientThread putRequests %s", request)
		ct.requestQueue.putRequest(request)
		time.Sleep(time.Duration(ct.random.Int()))
	}
}

type ServerThread struct {
	random       *rand.Rand
	requestQueue RequestQueue
}

func (st ServerThread) run() {
	for i := 0; i < 10000; i++ {
		request := st.requestQueue.getRequest()
		fmt.Println("ServerThread getRequests %s", request)
		time.Sleep(time.Duration(st.random.Int()))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	requestQueue := RequestQueue{queue: list.New()}
	clientThread := ClientThread{random: rand.New(rand.NewSource(time.Now().UnixNano())), requestQueue: requestQueue}
	serverThread := ServerThread{random: rand.New(rand.NewSource(time.Now().UnixNano())), requestQueue: requestQueue}
	go clientThread.run()
	go serverThread.run()

	select {}
}
