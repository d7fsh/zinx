package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/fatih/color"
	"zinx_demo/utils"
	"zinx_demo/ziface"
)

// 链接模块
type Connection struct {
	// 当前Conn隶属于那个Server
	TcpServer ziface.IServer
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 告知当前连接已经退出的/停止 channel(由reader告知Writer退出)
	ExitChan chan bool
	// 无缓冲的channel, 用于读写goroutine之间的消息通信
	msgChan chan []byte
	// 消息的管理MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler

	// 链接属性集合
	property map[string]interface{}
	// 保护连接属性的修改的锁
	propertyLock sync.RWMutex
}

// 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {

	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnManager().Add(c)

	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	color.Red("[Reader Goroutine is running...\n")
	defer fmt.Println("ConnID =", c.ConnID, " [Reader is exit], remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		dp := NewDataPack()

		// 读取客户端的Msg Head 二级制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("Read msg head error:", err)
			break
		}

		// 拆包, 得到 msgID和msgDataLen, 放置到msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}
		// 根据dataLen 再次读取Data, 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)
		// 得到当前conn数据的Request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制, 将消息发送给Worker工作池处理
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 从路由中, 找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID, 找到对应处理api业务, 执行
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// 写消息Goroutine, 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	color.Red("[Writer Goroutine is running...]\n")
	defer color.Yellow("%s conn Writer exit!\n", c.RemoteAddr().String())

	// 不断的阻塞的等待channel的消息, 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader已经退出, 此时Writer也要退出
			return
		}
	}
}

// 提供一个SendMsg方法, 将我们要发送给客户端的数据, 先进行封包, 再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed where send message!")
	}

	// 将data进行封包 msgDataLen|msgID|msgData
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Printf("Pack error msg id = %d\n", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID =", c.ConnID)
	// 启动从当前连接的读数据的业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建连接需要调用的处理业务, 执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID =", c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的 销毁连接之前需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭这个socket连接
	_ = c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从connManager中摘除掉
	c.TcpServer.GetConnManager().Remove(c)

	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

// 获取远程客户端的TCP状态, IP port
func (c *Connection) RemoteAddr() net.Addr {

	return c.Conn.RemoteAddr()
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 添加一个连接属性
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if v, ok := c.property[key]; ok {
		return v, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}
