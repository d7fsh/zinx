package znet

import (
	"errors"
	"fmt"
	"net"

	"github.com/fatih/color"
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
}

// 定义当前客户端连接的所绑定handleAPI(目前这个额handle是写死的, 以后优化应该由用户自定义)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 做一些回显功能
	fmt.Println("[Conn Handle] CallbackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("Write back buf err:", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// 启动服务器
func (s *Server) Start() {
	color.Cyan("[Start] Server Listener at [IP: %s, Port: %d] is starting\n", s.IP, s.Port)
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
			dealConn := NewConnection(conn, cid, CallBackToClient)
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

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
