package tcp_iface

type IMessageHandle interface {
	// DoMsgHandler 执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 委消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendMessageToTaskQueue 将消息交给TaskQueue，由worker进行处理
	SendMessageToTaskQueue(request IRequest)
}
