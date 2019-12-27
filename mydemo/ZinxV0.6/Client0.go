package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/fatih/color"
	"zinx_demo/znet"
)

/*
模拟客户端
*/
func main() {
	color.Cyan("client0 start...\n")
	time.Sleep(time.Second)

	// 1. 链接远程服务器, 得到一个conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		color.Red("client start err: %v", err)
		os.Exit(1)
	}

	for {
		// 发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.6 Client Test Message")))
		if err != nil {
			fmt.Println("client Pack error:", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 服务器应该回复message数据, MsgID:1 ping...消息
		// 1. 先读取流中的head部分, 得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}
		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 2. 再根据dataLen进行第二次读取, 将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}
			fmt.Printf("---------> recv Server msgId = %d, len = %d, data = %s", msg.Id, msg.DataLen, msg.Data)
		}

		time.Sleep(time.Second)
	}

}
