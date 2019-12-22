package main

import "zinx_demo/znet"

/*
基于Zinx框架开发的服务器端应用程序
*/
func main() {
	// 1. 创建一个Server句柄, 使用Zinx的api
	s := znet.NewServer("[zinx v0.2]")
	// 2,启动server
	s.Serve()
}
