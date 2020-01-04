package znet

import (
	"errors"
	"sync"

	"github.com/fatih/color"
	"github.com/d7fsh/zinx/ziface"
)

// 连接管理模块
type ConnManager struct {
	connections sync.Map // 管理的连接信息集合
	//connections map[uint32]ziface.IConnection // 管理的连接信息集合
	//connLock    sync.RWMutex                  // 保护连接集合的读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		//connections: make(map[uint32]ziface.IConnection),
		connections: sync.Map{},
	}
}

// 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map, 加写锁
	//cm.connLock.Lock()
	//defer cm.connLock.Unlock()

	// 将conn加入到connections中
	//cm.connections[conn.GetConnID()] = conn
	cm.connections.Store(conn.GetConnID(), conn)
	color.Yellow("connection [connID = %d] add to ConnManager successfully, conn num = %d\n", conn.GetConnID(), cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map, 加写锁
	//cm.connLock.Lock()
	//defer cm.connLock.Unlock()

	// 删除连接信息
	//delete(cm.connections, conn.GetConnID())
	cm.connections.Delete(conn.GetConnID())
	color.Yellow("connection [connID = %d] add to ConnManager successfully, conn num = %d\n", conn.GetConnID(), cm.Len())
}

// 根据connID获取链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源map, 加写锁
	//cm.connLock.Lock()
	//defer cm.connLock.Unlock()

	//if conn, ok := cm.connections[connID]; ok {
	if conn, ok := cm.connections.Load(connID); ok {
		return conn.(ziface.IConnection), nil
	} else {
		return nil, errors.New("connection [id =")
	}
}

// 得到当前连接总数
func (cm *ConnManager) Len() int {
	count := 0
	cm.connections.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// 删除并终止所有连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源map, 加写锁
	//cm.connLock.Lock()
	//defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	//for id, conn := range cm.connections {
	//	// 停止
	//	conn.Stop()
	//
	//	// 删除
	//	delete(cm.connections, id)
	//}
	// 删除conn并停止conn的工作
	cm.connections.Range(func(id, conn interface{}) bool {
		conn.(ziface.IConnection).Stop()
		cm.connections.Delete(id)
		return true
	})
	color.Yellow("clear all connections success! conn count = %d\n", cm.Len())
}
