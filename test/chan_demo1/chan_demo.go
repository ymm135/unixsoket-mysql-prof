package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup

func main() {

	msgQueue := make(chan string, 1)

	go listenQueue(msgQueue)

	num := 10
	waitGroup.Add(num)

	i := 0
	for i < num {
		msgQueue <- "data" + strconv.Itoa(i)
		i++
	}

	fmt.Println("queue len=", len(msgQueue), ",cap=", cap(msgQueue))
	waitGroup.Wait()
	fmt.Println(" == done == ")
}

func listenQueue(queue chan string) {
	fmt.Println("queue=", len(queue))
	for {
		select {
		case msg := <-queue:
			fmt.Println(msg)
			time.Sleep(time.Millisecond * 100)
			waitGroup.Done()
		}
	}
}
