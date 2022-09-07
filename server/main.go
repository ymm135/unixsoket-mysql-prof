package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unix-server/model"
	"unix-server/utils"

	"net/http"
	_ "net/http/pprof"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const socketFile = "/tmp/prof_sock"
const LookGid = true

var globalDb *gorm.DB
var msgQueue chan model.Employees

func main() {
	pid := os.Getpid()
	fmt.Println("-- unix socket server pid=", pid, " --")

	// 设置/配置
	fmt.Println(len(msgQueue))

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

	unixSocket := myutils.NewUnixSocket(socketFile, 1024)
	unixSocket.SetContextHandler(func(contexts string) string {
		paeseDataAndStore(contexts)
		return "ok"
	})

	go func() {
		unixSocket.StartServer()
	}()

	fmt.Println("-- unix socket server wait --")
	err = http.ListenAndServe("0.0.0.0:6060", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("-- unix socket server end --")
}

func paeseDataAndStore(context string) { // 多协程回调,每个回调都是一个协程 go this.HandleServerConn(c, string(data[0:nr]))
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
