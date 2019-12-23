package ziface

// IRouter
// 路由抽象接口
// 路由里的数据都是IRequest
type IRouter interface {
	// PreHandle 是在处理conn之前的钩子方法Hook
	PreHandle(request IRequest)
	// Handle 是处理conn业务的主方法hook
	Handle(request IRequest)
	// PostHandle 是处理conn业务之后的钩子方法hook
	PostHandle(request IRequest)
}
