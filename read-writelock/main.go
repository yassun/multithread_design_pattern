package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Data struct {
	buffer string
	lock   ReadWriteLock
}

func (d *Data) read() string {
	d.lock.readLock()
	buffer := d.doRead()
	d.lock.readUnlock()
	return buffer
}

func (d *Data) doRead() string {
	nebuf := d.buffer
	d.slowly()
	return nebuf
}

func (d *Data) write(c byte) {
	d.lock.writeLock()
	d.doWrite(c)
	d.lock.writeUnLock()
}

func (d *Data) doWrite(c byte) {
	size := len(d.buffer)
	d.buffer = ""
	for i := 0; i < size; i++ {
		d.buffer += string(c)
		d.slowly()
	}
}

func (d *Data) slowly() {
	time.Sleep(50 * time.Millisecond)
}

func NewData(size int) *Data {
	var s string
	for i := 0; i < size; i++ {
		s += "*"
	}

	readWriteLock := NewReadWriteLock()

	return &Data{buffer: s, lock: *readWriteLock}
}

type WriterThread struct {
	random *rand.Rand
	data   *Data
	filler string
	index  int
}

func (wt WriterThread) run() {
	for {
		c := wt.nextChar()
		wt.data.write(c)
		time.Sleep(time.Duration(wt.random.Int31()))
	}
}

func (wt *WriterThread) nextChar() byte {
	c := wt.filler[wt.index]
	wt.index++
	if wt.index >= len(wt.filler) {
		wt.index = 0
	}
	return c
}

type ReaderThread struct {
	data *Data
}

func (rt ReaderThread) run(name string) {
	for {
		fmt.Printf("%s reads %s\n", name, rt.data.read())
	}
}

type ReadWriteLock struct {
	readingReaders int
	waitingWriters int
	writingWriters int
	preferWriter   bool
	cond           *sync.Cond
}

func (rwl *ReadWriteLock) readLock() {
	rwl.cond.L.Lock()
	for rwl.writingWriters > 0 || (rwl.preferWriter && rwl.waitingWriters > 0) {
		rwl.cond.Wait()
	}
	rwl.readingReaders++
}

func (rwl *ReadWriteLock) readUnlock() {
	rwl.readingReaders--
	rwl.preferWriter = true
	rwl.cond.Signal()
	rwl.cond.L.Unlock()
}

func (rwl *ReadWriteLock) writeLock() {
	rwl.cond.L.Lock()
	rwl.waitingWriters++
	for rwl.readingReaders > 0 || rwl.writingWriters > 0 {
		rwl.cond.Wait()
	}

	rwl.waitingWriters--
	rwl.writingWriters++
}

func (rwl *ReadWriteLock) writeUnLock() {
	rwl.writingWriters--
	rwl.preferWriter = false

	rwl.cond.Signal()
	rwl.cond.L.Unlock()
}

func NewReadWriteLock() *ReadWriteLock {
	var mutex sync.Mutex
	return &ReadWriteLock{readingReaders: 0, writingWriters: 0, waitingWriters: 0, preferWriter: true, cond: sync.NewCond(&mutex)}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	data := NewData(10)
	go ReaderThread{data: data}.run("Thread-1")
	go ReaderThread{data: data}.run("Thread-2")
	go ReaderThread{data: data}.run("Thread-3")
	go ReaderThread{data: data}.run("Thread-4")
	go ReaderThread{data: data}.run("Thread-5")
	go ReaderThread{data: data}.run("Thread-6")

	go WriterThread{random: rand.New(rand.NewSource(time.Now().Unix())), data: data, filler: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", index: 0}.run()
	go WriterThread{random: rand.New(rand.NewSource(time.Now().Unix())), data: data, filler: "abcdefghijklmnopqrstuvwxyz", index: 0}.run()
	select {}

}
