package main

import (
	"log"
	"sync"
	"working/work"
	//"time"
)

type NamePrinter struct {
	name string
}

func (m *NamePrinter) Task() {
	log.Println(m.name)
	//time.Sleep(time.Second)
}

var names = []string{
	"steve",
	"bob",
	"mary",
	"therese",
	"jason",
	}

func main() {
	p := work.New(2)
	var wg sync.WaitGroup
	wg.Add(100*len(names))
	for i := 0; i < 100; i++ {
		for _, name := range names {
			np := NamePrinter{name:name}
			go func(){
				p.Run(&np)
				wg.Done()
			}()
		}
	}
	wg.Wait()
	p.Shutdown()
}