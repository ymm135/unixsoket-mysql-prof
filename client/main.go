package main

import (
	"fmt"
	"strconv"
	myutils "unix-server/utils"
)

const socketFile = "/tmp/prof_sock"

func main() {
	fmt.Println("--unix socket client --")
	unixSocket := myutils.NewUnixSocket(socketFile, 1024)
	count := 0
	max := 10000

	for count < max {
		unixSocket.ClientSendContext("1953-09-02,G" + strconv.Itoa(count) + ",Tester,M,1986-06-26")
		count++
	}
}
