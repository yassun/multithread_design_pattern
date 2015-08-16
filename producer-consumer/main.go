package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type MakerThread struct {
	random *rand.Rand
	table  *Table
	name   string
}

func (mt MakerThread) Run() {
	for {
		time.Sleep(time.Duration(mt.random.Int31()))
		cake := fmt.Sprintf("[ Cake No. %d by %s ]", nextId(), mt.name)
		mt.table.put(cake, mt.name)
	}
}

var id int = 0

func nextId() int {
	id = id + 1
	return id
}

type EaterThread struct {
	random *rand.Rand
	table  *Table
	name   string
}

func (et EaterThread) Run() {
	for {
		et.table.take(et.name)
		time.Sleep(time.Duration(et.random.Int31()))
	}
}

type Table struct {
	buffer []string
	tail   int
	head   int
	count  int
	cond   *sync.Cond
}

func (t *Table) put(cake string, name string) {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()

	fmt.Printf("%s puts %s\n", name, cake)
	for t.count >= len(t.buffer) {
		t.cond.Wait()
	}

	t.buffer[t.tail] = cake
	t.tail = (t.tail + 1) % len(t.buffer)
	t.count++
	t.cond.Signal()
}

func (t *Table) take(name string) string {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()

	for t.count <= 0 {
		t.cond.Wait()
	}

	cake := t.buffer[t.head]
	t.head = (t.head + 1) % len(t.buffer)
	t.count--

	fmt.Printf("%s takes %s\n", name, cake)
	t.cond.Signal()
	return cake
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	l := new(sync.Mutex)
	c := sync.NewCond(l)
	buffer := make([]string, 3)
	table := Table{buffer: buffer, head: 0, tail: 0, count: 0, cond: c}
	go MakerThread{random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "MakerThread-1"}.Run()
	go MakerThread{random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "MakerThread-2"}.Run()
	go MakerThread{random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "MakerThread-3"}.Run()
	go EaterThread{random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "EaterThread-1"}.Run()
	go EaterThread{random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "EaterThread-2"}.Run()
	go EaterThread{random: rand.New(rand.NewSource(time.Now().Unix())), table: &table, name: "EaterThread-3"}.Run()
	select {}
}
