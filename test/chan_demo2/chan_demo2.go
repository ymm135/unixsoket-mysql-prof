package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"
)

func main() {
	msgQueue := make(chan string, 1)
	go listenQueue(msgQueue, "queue1")
	go listenQueue(msgQueue, "queue2")

	num := 1000
	i := 0
	for i < num {
		go func() {
			msgQueue <- "data" + strconv.Itoa(i)
		}()
		i++
	}

	fmt.Println("queue len=", len(msgQueue), ",cap=", cap(msgQueue))
	err := http.ListenAndServe("0.0.0.0:6080", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(" == done == ")
}

func listenQueue(queue chan string, tag string) {
	fmt.Println("queue=", len(queue))
	for {
		select {
		case msg := <-queue:
			fmt.Println(tag, " data:", msg)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
