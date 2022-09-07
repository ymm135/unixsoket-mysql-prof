package main

import (
	"fmt"
	"time"
)

func main() {
	unix := time.Now().Unix()
	time.Sleep(time.Second)
	fmt.Println(time.Now().Unix() - unix)
}
