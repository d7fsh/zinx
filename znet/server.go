package znet

import (
	"fmt"
	"net"

	"github.com/fatih/color"
	"zinx_demo/utils"
	"zinx_demo/ziface"
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
}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
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

			// 将处理新连接业务方法和conn进行绑定, 得到我们的连接模块对象
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源, 状态或者一些已经开辟的链接信息, 进行停止回收
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
