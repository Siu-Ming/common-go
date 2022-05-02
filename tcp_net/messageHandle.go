package tcp_net

import (
	"fmt"
	"github.com/Siu-Ming/common-go/tcp_iface"
	"github.com/Siu-Ming/common-go/tcp_utils"
	"github.com/sirupsen/logrus"
	"strconv"
)

type MessageHandle struct {
	// 存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32]tcp_iface.IRouter
	// 业务工作Worker池的数量
	WorkerPoolSize uint32
	// Worker负责去任务消息队列
	TaskQueue []chan tcp_iface.IRequest
}

func NewMessageHandle() *MessageHandle {
	return &MessageHandle{
		Apis:           make(map[uint32]tcp_iface.IRouter),
		WorkerPoolSize: tcp_utils.GlobalObj.WorkerPoolSize,
		// 一个worker对应一个queue
		TaskQueue: make([]chan tcp_iface.IRequest, tcp_utils.GlobalObj.WorkerPoolSize),
	}
}

// DoMsgHandler 马上以非阻塞方式处理消息
func (m *MessageHandle) DoMsgHandler(request tcp_iface.IRequest) {
	handler, ok := m.Apis[request.GetMsgId()]
	if !ok {
		logrus.Error("api msgId=", request.GetMsgId(), "没有创建!")
		return
	}
	//执行对应处理方法
	handler.BeforeHandle(request)
	handler.Handle(request)
	handler.AfterHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (m *MessageHandle) AddRouter(msgId uint32, router tcp_iface.IRouter) {
	// 1. 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api, msgId=" + strconv.Itoa(int(msgId)))
	}
	// 添加msg与API的绑定关系
	m.Apis[msgId] = router
	fmt.Println("Add api msgId", msgId)
}

// StartWorkerPool 启动worker工作池
func (m *MessageHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan tcp_iface.IRequest, tcp_utils.GlobalObj.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞等待对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// SendMessageToTaskQueue 将消息交给TaskQueue，由worker进行处理
func (m *MessageHandle) SendMessageToTaskQueue(request tcp_iface.IRequest) {

	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	// 得到需要处理此条链接的workerID
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("ADD ConnID=", request.GetConnection().GetConnID(), "request msgID=", request.GetMsgId(), "to workerID=", workerID)
	// 将请求消息发送给任务队列
	m.TaskQueue[workerID] <- request
}

// StartOneWorker 启动一个Worker工作流程
func (m *MessageHandle) StartOneWorker(workerID int, taskQueue chan tcp_iface.IRequest) {
	fmt.Println("Worker ID = ", workerID, "启动")
	for true {
		select {
		case requests := <-taskQueue:
			m.DoMsgHandler(requests)
		}
	}
}
