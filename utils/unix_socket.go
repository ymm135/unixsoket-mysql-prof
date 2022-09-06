package socket

import (
	"fmt"
	"net"
	"os"
	"time"
)

type UnixSocket struct {
	filename string
	bufsize  int
	handler  func(string) string
}

func NewUnixSocket(filename string, size ...int) *UnixSocket {
	size1 := 10480
	if size != nil {
		size1 = size[0]
	}
	us := UnixSocket{filename: filename, bufsize: size1}
	return &us
}

func (this *UnixSocket) createServer() {
	fmt.Println("socket监听执行========================================")
	os.Remove(this.filename)
	addr, err := net.ResolveUnixAddr("unixgram", this.filename)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	c, err := net.ListenUnixgram("unixgram", addr)
	defer c.Close()
	if err != nil {
		panic("Cannot listen to unix domain socket: " + err.Error())
	}
	os.Chmod(this.filename, 0666)
	for {
		data := make([]byte, 4096)
		nr, _, err := c.ReadFrom(data)
		if err != nil {
			fmt.Printf("conn.ReadFrom error: %s\n", err)
			return
		}
		go this.HandleServerConn(c, string(data[0:nr]))
	}

}

//接收连接并处理
func (this *UnixSocket) HandleServerConn(conn net.Conn, data string) {
	this.HandleServerContext(data)
}

func (this *UnixSocket) SetContextHandler(f func(string) string) {
	this.handler = f
}

//接收内容并返回结果
func (this *UnixSocket) HandleServerContext(context string) string {
	if this.handler != nil {
		return this.handler(context)
	}
	now := time.Now().String()
	return now
}

func (this *UnixSocket) StartServer() {
	this.createServer()
}

//客户端
func (this *UnixSocket) ClientSendContext(context string) {
	addr, err := net.ResolveUnixAddr("unixgram", this.filename)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	//拔号
	c, err := net.DialUnix("unixgram", nil, addr)
	if err != nil {
		panic("DialUnix failed.")
	}
	//写出
	_, err = c.Write([]byte(context))
	if err != nil {
		panic("Writes failed.")
	}
}
