package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Data struct {
	fileName string
	content  string
	changed  bool
	mutex    sync.Mutex
}

func (d *Data) Change(newContent string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.content = newContent
	d.changed = true
}

func (d *Data) Save(callsName string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.changed {
		return
	}
	d.doSave(callsName)
	d.changed = false
}

func (d *Data) doSave(callsName string) {
	fmt.Println("calls：" + callsName + "doSave, fileName：" + d.fileName + ", content：" + d.content)
}

type SaverThread struct {
	data *Data
	name string
}

func (st SaverThread) Run() {
	for {
		st.data.Save(st.name)
		time.Sleep(1 * time.Second)
	}
}

type ChangerThread struct {
	data   *Data
	random *rand.Rand
	name   string
}

func (ct ChangerThread) Run() {
	for i := 0; ; i++ {
		ct.data.Change(fmt.Sprintf("No. %d", i))

		// 仕事のつもり
		time.Sleep(time.Duration(ct.random.Int31()))

		ct.data.Save(ct.name)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	data := &Data{fileName: "data.txt", content: "(empty)", changed: true, mutex: *new(sync.Mutex)}
	changerThread := ChangerThread{data: data, random: rand.New(rand.NewSource(time.Now().Unix())), name: "ChangerThread"}
	saverThread := SaverThread{data: data, name: "SaverThread"}

	go changerThread.Run()
	go saverThread.Run()

	select {}
}
