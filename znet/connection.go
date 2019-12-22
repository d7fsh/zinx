package znet

import (
	"net"

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
	// 当前连接锁绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	// 告知当前连接已经退出的/停止 channel
	ExitChan chan bool
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackApi,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
	return c
}
