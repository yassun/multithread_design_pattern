package main

import (
	"fmt"
	"time"
)

type TsLog struct{}

func (tl TsLog) printLn(s string) {
	fmt.Println(s)
}

func (tl TsLog) close() {
	fmt.Println("==== End of log ====")
}

type Log struct {
	tsLogCollection ThredLocal
}

func (l Log) printLn(s string) {
	tsLog := l.getTslog()
	tsLog.printLn(s)
}

func (l Log) close() {
	tsLog := l.getTslog()
	tsLog.printLn("==== End of log ====")
}

func (l Log) getTslog() *TsLog {
	tsLog := l.tsLogCollection.get()

	if tsLog == nil {
		tsLog := TsLog{}
		l.tsLogCollection.set(tsLog)
	}

	return tsLog
}

// TODO
type ThredLocal struct {
}

func (tl ThredLocal) get() *TsLog {
	return &TsLog{}
}

func (tl *ThredLocal) set(tslog TsLog) {
}

type ClientThread struct {
	name string
}

func (ct ClientThread) run() {
	fmt.Printf("%s BEGIN", ct.name)

	for i := 0; i < 10; i++ {
		//Log.printLn(fmt.Sprintf("i = %d", i))
		time.Sleep(100 * time.Millisecond)
	}

	//Log.close()
	fmt.Printf("%s END", ct.name)
}
