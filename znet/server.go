package znet

// IServer接口实现
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的IP地址
	IP string
	// 服务器监听的端口
	Port int
}

// 启动服务器
func (s *Server) Start() {

}

// 停止服务器
func (s *Server) Stop() {

}

// 运行服务器
func (s *Server) Serve() {

}
