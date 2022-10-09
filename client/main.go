package main

import (
	"fmt"
	"sync"
	"time"
	myutils "unix-server/utils"
)

const socketFile = "/data/socket/obtain_proto_sock1"

var wg sync.WaitGroup

func main() {
	fmt.Println("--unix socket client --")
	unixSocket := myutils.NewUnixSocket(socketFile, 1024)
	count := 0
	max := 1
	wg.Add(max)

	for count < max {
		go func() {
			go unixSocket.ClientSendContext("'2000-01-01 20:54:43','2000-01-01 20:54:43','c8:5b:76:3e:a5:5d','3e:c0:18:af:9d:d9', 4,'192.168.99.45',3812,'192.168.99.31',502,6,'TCP',64,'','modbus'----func----''----")
			wg.Done()
		}()
		time.Sleep(time.Millisecond * 1)
		count++
	}
	fmt.Println("clint wait")
	wg.Wait()
	fmt.Println("clint done")
}
