package znet

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
)

// 只是负责测试dataPack拆包, 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	// 1. 创建socketTCP server
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		log.Panic("server listen err:", err)
	}

	// 创建一个 go 承载从客户端处理业务
	go func() {
		// 2. 从客户端读取数据, 拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error:", err)
				return
			}

			go func(c net.Conn) {
				// 处理客户端的请求
				//  -----> 拆包的过程 <------
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 1. 第一次从conn读, 把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error:", err)
						return
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack error:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// Msg是有数据的. 需要进行第二次读取
						// 2. 第二次从conn读, 根据head中的dataLen在读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						// 根据dataLen的长度再次从io流中读取
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("server unpack data error:", err)
							return
						}
						//完整的消息已经读取完毕
						fmt.Printf("-------> Recv MsgID: %d, dataLen = %d, data = %s\n", msg.Id, msg.DataLen, msg.Data)
					}
				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}
	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟粘包过程, 封装两个msg一起发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}
	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'g', 'o', 'w', 'o', 'r', 'l', 'd'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error:", err)
		return
	}
	// 将两个包黏在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	_, _ = conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
