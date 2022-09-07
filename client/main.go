package main

import (
	"fmt"
	"strconv"
	"sync"
	myutils "unix-server/utils"
)

const socketFile = "/tmp/prof_sock"

var wg sync.WaitGroup

func main() {
	fmt.Println("--unix socket client --")
	unixSocket := myutils.NewUnixSocket(socketFile, 1024)
	count := 0
	max := 100
	wg.Add(max)

	for count < max {
		unixSocket.ClientSendContext("1953-09-02,G" + strconv.Itoa(count) + ",Tester,M,1986-06-26")
		count++
		wg.Done()
	}
	fmt.Println("clint wait")
	wg.Wait()
	fmt.Println("clint done")
}
