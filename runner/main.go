package main

import (
	"log"
	"time"
	"runner"
)

func main() {
	log.Println("Starting work.")
	
	r = run
}

func createTask() func(int) {
	return func(id int) {
		log.Println("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id)*time.Second)
	}
}