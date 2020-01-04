package main

import (
	"fmt"

	"github.com/fatih/color"
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

// 创建连接之后执行的钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	color.Green("-----> DoConnectionBegin is Called...\n")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin.....")); err != nil {
		color.Red("%v\n", err)
	}

	// 给当前的连接设置一些属性
	color.Yellow("Set conn Name, ....\n")
	conn.SetProperty("Name", "..........b2.pub........")
	conn.SetProperty("Home", "http://www.b2.pub")
}

// 断开连接之前需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	color.Yellow("-----> DoConnectionLost is Called...\n")
	color.Yellow("conn ID = %d\n", conn.GetConnID())

	// 获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		color.Yellow("Name = %s\n", name)
	}
}

func main() {
	// 1. 创建一个server句柄, 使用zinx的api
	s := znet.NewServer()

	// 2. 注册连接的hookFunc
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3. 给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.AddRouter(2, &HelloRouter{})
	s.AddRouter(3, &HelloRouter{})

	// 4. 启动server
	s.Serve()
}
