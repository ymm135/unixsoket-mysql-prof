package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	myutils "unix-server/utils"
)

const socketFile = "/tmp/prof_sock"

var wg sync.WaitGroup

func main() {
	fmt.Println("--unix socket client --")
	unixSocket := myutils.NewUnixSocket(socketFile, 1024*1024)
	count := 0
	max := 2000000
	wg.Add(max)

	for count < max {
		go func() {
			go unixSocket.ClientSendContext("1953-09-02,G" + strconv.Itoa(count) + ",Tester,M,1986-06-26")
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 1)
		count++
	}
	fmt.Println("clint wait")
	wg.Wait()
	fmt.Println("clint done")
}
