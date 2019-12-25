package ziface

/*
封包, 拆包 模块
直接面向TCP连接中的数据流, 用于处理解决TCP粘包问题
*/
type IDataPacker interface {
	// 获取包的头的长度的方法
	GetHeadLen() uint32
	// 封包方法
	Pack(IMessage) ([]byte, error)
	// 拆包方法
	Unpack([]byte) (IMessage, error)
}
