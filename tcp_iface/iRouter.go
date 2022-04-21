package tcp_iface

type IRouter interface {
	// BeforeHandle 在conn业务之前的钩子方法
	BeforeHandle(request IRequest)
	// Handle 处理conn业务的方法
	Handle(request IRequest)
	// AfterHandle 处理conn业务之后的钩子方法
	AfterHandle(request IRequest)
}
