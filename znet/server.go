package znet

import (
	"fmt"
	"net"

	"github.com/d7fsh/zinx/utils"
	"github.com/d7fsh/zinx/ziface"
	"github.com/fatih/color"
)

// IServer接口实现
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的IP地址
	IP string
	// 服务器监听的端口
	Port int
	// 当前server消息管理模块, 用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 该server的连接管理器
	ConnMan ziface.IConnManager
	// 该Server创建链接之后自动调用Hook函数 --- OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 该server销毁之前自动调用Hook函数 -- OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

/*
初始化Server模块的方法
*/
func NewServer() ziface.IServer {

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMan:    NewConnManager(),
	}
	return s
}

// 启动服务器
func (s *Server) Start() {
	color.Cyan("[Zinx] ServerName: %s, listener at IP: %s, port: %d is starting",
		s.Name,
		s.IP,
		s.Port,
	)
	color.Cyan("[Zinx] Version %s, MaxConn: %d, MaxPackageSize: %d", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0. 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			color.Red("resolve tcp addr error:", err)
			return
		}
		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			color.Red("listen ", s.IPVersion, " error", err)
			return
		}
		color.Cyan("Start Zinx server %s success, Listening...\n", s.Name)
		var cid uint32
		cid = 0

		// 3. 阻塞等待客户端连接, 处理客户端连接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				color.Red("Accept error:", err)
				continue
			}

			// 设置最大连接个数的判断, 如果超过最大连接的数量, 那么则关闭此连接
			if s.ConnMan.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个错误超出最大连接的错误包
				color.Red("Too Many Connections Max Conn = %d\n", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 将处理新连接业务方法和conn进行绑定, 得到我们的连接模块对象
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源, 状态或者一些已经开辟的链接信息, 进行停止回收
	color.Red("[STOP] Zinx server name %s\n", s.Name)
	s.ConnMan.ClearConn()
}

// 运行服务器
func (s *Server) Serve() {
	// 启动服务器
	s.Start()

	//TODO 做一些启动之后的额外业务

	// 阻塞状态
	select {}
}

// 路由功能: 给当前的服务注册一个路由方法, 供客户端的连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success.")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMan
}

// 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("------> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---------> Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
