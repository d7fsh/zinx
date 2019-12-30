package main

import (
	"fmt"

	"zinx_demo/ziface"
	"zinx_demo/znet"
)

/*
	基于Zinx框架来开发的服务端应用程序
*/

// ping test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call ----Ping---- Router Handle")
	// 先读取客户端的数据, 在回写ping...ping...ping...
	fmt.Printf("recv from client: msgId = %d, data = %s\n", request.GetMsgID(), request.GetData())

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping...\r\n"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

// Hello Router
func (p *HelloRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call ----Hello---- Router Handle")
	// 先读取客户端的数据, 在回写ping...ping...ping...
	fmt.Printf("recv from client: msgId = %d, data = %s\n", req.GetMsgID(), req.GetData())

	err := req.GetConnection().SendMsg(201, []byte("hello hello hello\r\n"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1. 创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[ZinxV0.8]")
	// 2. 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 2, 启动server
	s.Serve()
}
