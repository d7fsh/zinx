package znet

import (
	"fmt"
	"net"

	"zinx_demo/utils"
	"zinx_demo/ziface"
)

// 链接模块
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 告知当前连接已经退出的/停止 channel
	ExitChan chan bool
	// 该链接处理的方法Router
	Router ziface.IRouter
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("ConnID =", c.ConnID, " Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中, 最大512 字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv vuf err", err)
			continue
		}

		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			data: buf,
		}

		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			// 从路由中找到注册绑定的Conn对应的router
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(req)

	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID =", c.ConnID)
	// 启动从当前连接的读数据的业务
	go c.StartReader()

	// TODO 启动从当前连接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID =", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭这个socket连接
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn
}

func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
