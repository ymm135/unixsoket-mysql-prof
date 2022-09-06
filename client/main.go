package main

import (
	"fmt"
	"strconv"
	socket "unix-server/utils"
)

const socketFile = "/tmp/prof_sock"

func main() {
	fmt.Println("--unix socket client --")
	unixSocket := socket.NewUnixSocket(socketFile, 1024)
	count := 0

	for count < 1 {
		unixSocket.ClientSendContext("1953-09-02,G" + strconv.Itoa(count) + ",Tester,M,1986-06-26")
		count++
	}
}
