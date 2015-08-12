package main

import (
	"fmt"
	"sync"
	"time"
)

type Gate struct {
	counter int
	name    string
	address string
	m       *sync.Mutex
}

func (g *Gate) String() string {
	return fmt.Sprintf("No.%d: %s, %s", g.counter, g.name, g.address)
}

func (g *Gate) check() {
	if g.name[0] != g.address[0] {
		println("***** BROKEN ***** " + g.String())
	} else {
		println("***** NOT BROKEN ***** " + g.String())
	}
}

func (g *Gate) Pass(name, address string) {
	g.m.Lock()
	defer g.m.Unlock()

	g.counter++
	g.name = name
	g.address = address
	g.check()
}

type UserThread struct {
	gate      *Gate
	myName    string
	myAddress string
}

func (thread UserThread) Start() {
	println(fmt.Sprintf("%s BEGIN", thread.myName))
	for {
		time.Sleep(100 * time.Millisecond)
		thread.gate.Pass(thread.myName, thread.myAddress)
	}
}

func main() {
	println("CTRL+C to exit.")
	m := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	gate := Gate{counter: 0, name: "Nobody", address: "Nowhere", m: m}
	alice := UserThread{gate: &gate, myName: "Alice", myAddress: "Alaska"}
	bobby := UserThread{gate: &gate, myName: "Bobby", myAddress: "Brazil"}
	chris := UserThread{gate: &gate, myName: "Chris", myAddress: "Canada"}

	wg.Add(1)
	go alice.Start()

	wg.Add(1)
	go bobby.Start()

	wg.Add(1)
	go chris.Start()

	wg.Wait()

}
