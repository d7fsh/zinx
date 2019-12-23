package ziface

/*
	将请求的消息封装到一个Message中, 定义一个抽象的接口
*/
type IMessage interface {
	// 获取消息的Id
	GetMsgId() uint32
	// 获取消息的长度
	GetMsgLen() uint32
	// 获取消息的内容
	GetData() []byte

	// 设置消息的Id
	SetMsgId(uint32)
	// 设置消息的内容
	SetMsgLen(uint32)
	// 设置消息的长度
	SetData([]byte)
}
