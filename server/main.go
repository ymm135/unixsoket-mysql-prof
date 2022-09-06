package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unix-server/model"
	socket "unix-server/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const socketFile = "/tmp/prof_sock"

var wg sync.WaitGroup //定义一个同步等待的组
var globalDb *gorm.DB

func main() {
	fmt.Println("-- unix socket server --")

	// 连接mysql
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:root@tcp(10.25.10.125:3306)/test_prof?charset=utf8mb4&parseTime=True&loc=Local"
	dialector := mysql.New(mysql.Config{
		DSN: dsn,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println("connect mysql error!")
		return
	}
	globalDb = db
	fmt.Println(globalDb.Name())

	unixSocket := socket.NewUnixSocket(socketFile, 1024)
	unixSocket.SetContextHandler(func(contexts string) string {
		paeseDataAndStore(contexts)
		return "ok"
	})

	wg.Add(1)
	go func() {
		unixSocket.StartServer()
		wg.Done()
	}()

	fmt.Println("-- unix socket server wait --")
	wg.Wait()
	fmt.Println("-- unix socket server end --")
}

func paeseDataAndStore(context string) {
	fmt.Println("recvData:", context)
	fields := strings.Split(context, ",")

	birthDate, _ := time.ParseInLocation("2006-01-02", fields[0], time.Local)
	hireDate, _ := time.ParseInLocation("2006-01-02", fields[4], time.Local)

	employee := model.Employees{
		Id:        0,
		BirthDate: birthDate,
		FirstName: fields[1],
		LastName:  fields[2],
		//Gender:    []rune(fields[3])[0],
		Gender:   fields[3],
		HireDate: hireDate,
	}

	err := globalDb.Table("employees").Create(&employee).Error
	if err != nil {
		fmt.Println(err.Error())
	}
}