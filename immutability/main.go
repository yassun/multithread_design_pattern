package main

import (
	"fmt"
	"runtime"
)

type Person struct {
	name, address string
}

func (p Person) String() string {
	return fmt.Sprintf("[ Person: name = %s, address = %s ]", p.name, p.address)
}

type PrintPersonThread struct {
	person Person
}

func (t PrintPersonThread) Run() {
	for {
		fmt.Println(fmt.Sprintf("prints %s", t.person))
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	alice := Person{name: "Alice", address: "Alaska"}

	go PrintPersonThread{person: alice}.Run()
	go PrintPersonThread{person: alice}.Run()
	go PrintPersonThread{person: alice}.Run()

	select {}
}
