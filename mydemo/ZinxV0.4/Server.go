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

// Test PreHandle
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("[Before] ping...\n"))
	if err != nil {
		fmt.Println("call back [before] ping error:", err)
	}
}

// Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call back [ping] error:", err)
	}
}

// Test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("[After] ping...\n"))
	if err != nil {
		fmt.Println("call back After ping error:", err)
	}
}

func main() {
	// 1. 创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[ZinxV0.3]")
	// 2. 给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 2, 启动server
	s.Serve()
}
