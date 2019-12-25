package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"zinx_demo/utils"
	"zinx_demo/ziface"
)

/*
封包, 拆包的具体模块
封包的格式: 前四位存储数据长度, 接着四位存储数据包Id, 最后 5位存储数据
前 8 位组成head, 后 5 位组成body
*/
type DataPack struct{}

// 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头的长度的方法
func (d *DataPack) GetHeadLen() uint32 {
	// dataLen uint32 4字节 + id uint32 4字节
	return 8
}

// 封包方法
// |dataLen|msgID|data|
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	// 将dataLen写入dataBuf 中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将MsgId写入dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将Data数据写入 dataBuf中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// 拆包方法
// 将包的Head信息读取出来, 之后在根据head信息中的data的长度, 再进行一次读
func (d *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	reader := bytes.NewReader(data)

	// 只解压head信息, 得到dataLen和MsgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}
