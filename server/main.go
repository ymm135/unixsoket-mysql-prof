package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unix-server/model"
	myutils "unix-server/utils"

	"net/http"
	_ "net/http/pprof"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const socketFile = "/tmp/prof_sock"

var globalDb *gorm.DB
var msgQueue chan string

var lastSecondTime int64
var currHandlerCount int64
var lastHandlerCount int64

var saveDataQueue []model.Employees

var wg sync.WaitGroup

func main() {
	pid := os.Getpid()
	fmt.Println("-- unix socket server pid=", pid, " --")

	// 设置/配置
	lastSecondTime = time.Now().Unix()
	msgQueue = make(chan string, 10000)
	saveDataQueue = make([]model.Employees, 0)
	fmt.Println(len(msgQueue))

	// 连接mysql
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:root@tcp(10.25.10.125:3306)/test_prof?charset=utf8mb4&parseTime=True&loc=Local"
	dialector := mysql.New(mysql.Config{
		DSN: dsn,
	})

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:      newLogger,
		PrepareStmt: true, //
	})
	if err != nil {
		fmt.Println("connect mysql error!")
		return
	}
	//globalDb = db
	globalDb = db.Session(&gorm.Session{PrepareStmt: true})
	fmt.Println(globalDb.Name())

	unixSocket := myutils.NewUnixSocket(socketFile, 1024*1024)
	unixSocket.SetContextHandler(func(contexts string) string {
		// 多协程 入队
		msgQueue <- contexts
		return "ok"
	})

	go unixSocket.StartServer()
	go handleDataLoop()

	fmt.Println("-- ListenAndServe --")
	err = http.ListenAndServe("0.0.0.0:6060", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("-- unix socket server end --")
}

func handleDataLoop() {
	for {
		select {
		case data := <-msgQueue:
			paeseDataAndStore(data)
		}
	}
}

func paeseDataAndStore(context string) { // 多协程回调,每个回调都是一个协程 go this.HandleServerConn(c, string(data[0:nr]))
	//fmt.Println("recvData:", context)
	now := time.Now().Unix()
	if now-lastSecondTime >= 1 {
		fmt.Println("paeseDataAndStore handler data", currHandlerCount-lastHandlerCount, "pps", time.Now().String())
		lastSecondTime = now
		lastHandlerCount = currHandlerCount
	}

	currHandlerCount++

	fields := strings.Split(context, ",")
	//birthDate, _ := time.ParseInLocation("2006-01-02", fields[0], time.Local)
	hireDate, _ := time.ParseInLocation("2006-01-02", fields[4], time.Local)

	employee := model.Employees{
		Id:        0,
		BirthDate: time.Now(),
		FirstName: fields[1],
		LastName:  fields[2],
		Gender:    fields[3],
		HireDate:  hireDate,
	}

	saveDataQueue = append(saveDataQueue, employee)

	// 单个处理
	//err := globalDb.Table("employees").Create(&employee).Debug().Error
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	routineNum := 20
	if len(saveDataQueue) >= 10000 {
		wg.Add(routineNum)
		count := 1
		startInsertTime := time.Now()

		for count <= routineNum {
			go batchInsertData(count)
			count++
		}
		// 等所有数据写入完成
		wg.Wait()

		gap := time.Now().Unix() - startInsertTime.Unix()
		fmt.Println("batch insert data,cost ", gap, "s,avg", (float64)(gap)/(float64)(routineNum), "s")
		// 清空数据
		saveDataQueue = saveDataQueue[:0]

	}

}

func batchInsertData(index int) {
	currTime := time.Now().Unix()
	fmt.Println(index, "save data", len(saveDataQueue))
	err := globalDb.Table("employees" + strconv.Itoa(index)).Create(&saveDataQueue).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(index, "save data end 耗时:", time.Now().Unix()-currTime, "s")

	wg.Done()
}
