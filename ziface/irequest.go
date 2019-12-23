package ziface

/*
IRequest 接口
把客户端请求的连接信息, 和请求数据包装到了一个Request中
*/
type IRequest interface {
	// 得到当前连接
	GetConnection() IConnection
	// 得到请求的消息数据
	GetData() []byte
}
