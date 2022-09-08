package main

import (
	"fmt"
	"github.com/Jeffail/tunny"
	"sync"
	"time"
	myutils "unix-server/utils"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("main gid", myutils.GetGID())
	pool := tunny.NewFunc(2, func(in interface{}) interface{} {
		fmt.Println("gid", myutils.GetGID())
		intVal := in.(int)
		time.Sleep(time.Millisecond * 1000)
		wg.Done()
		return intVal * 2
	})
	defer pool.Close()

	count := 0
	num := 100
	wg.Add(num)
	for count < num {
		go pool.Process(10)
		//fmt.Println(ret)
		count++
	}

	wg.Wait()
}
