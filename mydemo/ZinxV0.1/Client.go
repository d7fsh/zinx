package main

import (
	"net"
	"os"
	"time"

	"github.com/fatih/color"
)

/*
模拟客户端
*/
func main() {
	color.Cyan("client start...\n")
	time.Sleep(time.Second)
	// 1. 链接远程服务器, 得到一个conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		color.Red("client start err: %v", err)
		os.Exit(1)
	}
	for {
		// 2. 调用Write写数据
		_, err = conn.Write([]byte("Hello zinx v0.1.."))
		if err != nil {
			color.Red("write conn err: %v\n", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			color.Red("read buf error: %v", err)
			return
		}

		color.Cyan("server cal back: %s, cnt = %d\n", buf, cnt)
		time.Sleep(time.Second)
	}

}
