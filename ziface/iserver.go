package ziface

// 定义一个服务器接口

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 路由功能, 给当前的服务注册一个路由方法, 工客户端的连接处理使用
	AddRouter(msgID uint32, router IRouter)
	// 获取当前连接的管理器
	GetConnManager() IConnManager
	// 注册OnConnStart钩子函数的方法
	SetOnConnStart(func(conn IConnection))
	// 注册OnConnStop钩子函数的方法
	SetOnConnStop(func(conn IConnection))
	// 调用OnConnStart钩子函数的方法
	CallOnConnStart(conn IConnection)
	// 调用OnConnStart钩子函数的方法
	CallOnConnStop(conn IConnection)
}
