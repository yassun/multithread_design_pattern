package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const WorkerCount int = 5

type ClientThread struct {
	name    string
	channel chan Request
	random  *rand.Rand
}

func (ct *ClientThread) Run() {
	for i := 0; ; i++ {
		request := NewRequest(ct.name, i)
		ct.channel <- *request
		time.Sleep(time.Duration(ct.random.Int31()))
	}
}

func NewClientThread(name string, channel chan Request) *ClientThread {
	return &ClientThread{name: name, channel: channel, random: rand.New(rand.NewSource(time.Now().Unix()))}
}

type Request struct {
	name   string
	number int
	random *rand.Rand
}

func NewRequest(name string, number int) *Request {
	return &Request{name: name, number: number, random: rand.New(rand.NewSource(time.Now().Unix()))}
}

func (r *Request) Execute() {
	fmt.Println(" executes ", r)
	time.Sleep(time.Duration(r.random.Int31()))
}

func (r *Request) String() string {
	return fmt.Sprintf("[ Request from %s No. %d ]", r.name, r.number)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	channel := make(chan Request, 100)
	go NewClientThread("Alice", channel).Run()
	go NewClientThread("Bobby", channel).Run()
	go NewClientThread("Chris", channel).Run()

	for i := 0; i < WorkerCount; i++ {
		go func(ch <-chan Request) {
			for {
				request := <-ch
				request.Execute()
			}
		}(channel)
	}

	select {}

}
