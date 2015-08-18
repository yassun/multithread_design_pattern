package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Host struct {
	helper Helper
}

func (h Host) request(count int, c byte, wg *sync.WaitGroup) {
	fmt.Printf("    request(%d, %c) BEGIN\n", count, c)
	go func() {
		h.helper.handle(count, c)
		wg.Done()
	}()
	fmt.Printf("    request(%d, %c) END\n", count, c)
}

type Helper struct{}

func (h Helper) handle(count int, c byte) {
	fmt.Printf("        handle(%d, %c) BEGIN\n", count, c)
	for i := 0; i < count; i++ {
		h.slowly()
		fmt.Printf("%c", c)
	}
	fmt.Println("")
	fmt.Printf("        handle(%d, %c) END\n", count, c)
}

func (h Helper) slowly() {
	time.Sleep(100 * time.Millisecond)
}

func newHost() *Host {
	return &Host{helper: Helper{}}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("main BEGIN")

	wg := new(sync.WaitGroup)
	host := newHost()

	wg.Add(1)
	host.request(10, 'A', wg)

	wg.Add(1)
	host.request(20, 'B', wg)

	wg.Add(1)
	host.request(30, 'C', wg)

	wg.Wait()

}
