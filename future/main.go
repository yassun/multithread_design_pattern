package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Host struct{}

func (h *Host) request(count int, c byte) Data {
	fmt.Printf("    request(%d, %s) BEGIN\n", count, string(c))

	future := NewFutureData()
	go func() {
		var realData RealData
		realData.makeRealData(count, c)
		future.setRealData(realData)
	}()

	fmt.Printf("    request(%d, %s) END\n", count, string(c))
	return &future
}

type Data interface {
	getContent() string
}

type FutureData struct {
	realData RealData
	ready    bool
	cond     *sync.Cond
}

func (fd *FutureData) setRealData(realData RealData) {
	fd.cond.L.Lock()
	defer fd.cond.L.Unlock()

	if fd.ready {
		return
	}

	fd.realData = realData
	fd.ready = true

	fd.cond.Signal()
}

func (fd *FutureData) getContent() string {
	fd.cond.L.Lock()
	defer fd.cond.L.Unlock()

	for !fd.ready {
		fd.cond.Wait()
	}

	return fd.realData.getContent()
}

func NewFutureData() FutureData {
	var mutex sync.Mutex
	return FutureData{cond: sync.NewCond(&mutex), ready: false}
}

type RealData struct {
	content string
}

func (rd *RealData) makeRealData(count int, c byte) {
	fmt.Printf("        making RealData(%d, %s) BEGIN\n", count, string(c))
	buffer := make([]byte, count)
	for i := 0; i < count; i++ {
		buffer[i] = c
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("        making RealData(%d, %s) END\n", count, string(c))
	rd.content = string(buffer)
}

func (rd *RealData) getContent() string {
	return rd.content
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("main BEGIN")
	var host Host
	data1 := host.request(10, 'A')
	data2 := host.request(20, 'B')
	data3 := host.request(30, 'C')

	fmt.Println("main otherJob BEGIN")
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("main otherJob END")

	fmt.Printf("data1 = %s\n", data1.getContent())
	fmt.Printf("data2 = %s\n", data2.getContent())
	fmt.Printf("data3 = %s\n", data3.getContent())

	fmt.Println("main END")

}
