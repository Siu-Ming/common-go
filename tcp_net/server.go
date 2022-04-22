package tcp_net

import (
	"fmt"
	"github.com/Siu-Ming/common-go/tcp_iface"
	"github.com/Siu-Ming/common-go/tcp_utils"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	// 服务器名称
	Name string
	// IP
	IP string
	// 端口
	Port int
	// IP版本
	IPVersion string
	// Router    tcp_iface.IRouter
	// 绑定msgid和对应的吹业务API关系
	MessageHandle tcp_iface.IMessageHandle
	// 当前server链接管理器
	ConnMgr tcp_iface.IConnectionManager

	//该Server的连接创建时Hook函数
	OnConnStart func(connection tcp_iface.IConnection)
	// 该Server的连接断开时的Hook函数
	OnConnStop func(connection tcp_iface.IConnection)
}

func NewServer(name string) *Server {
	server := &Server{
		Name:          name,
		IP:            tcp_utils.GlobalObj.Host,
		Port:          tcp_utils.GlobalObj.Port,
		IPVersion:     "tcp4",
		MessageHandle: NewMessageHandle(),
		ConnMgr:       NewConnectionManager(),
	}
	return server
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("服务启动 Name: %s,监听IP:%s,服务端口 %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[TCP-framework] 版本:%s,最大连接数:%d,最大数据包长度:%d\n",
		tcp_utils.GlobalObj.Version,
		tcp_utils.GlobalObj.MaxConn,
		tcp_utils.GlobalObj.MaxPacketSize)

	// 0.启动worker工作池机制
	s.MessageHandle.StartWorkerPool()

	// 1. 获取TCP的Address
	tcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		logrus.Error("resolve tcp addr error", err)
		return
	}

	// 2. 监听服务器地址
	listenTCP, err := net.ListenTCP(s.IPVersion, tcpAddr)
	if err != nil {
		logrus.Error("listen", s.IPVersion, " tcp error", err)
		return
	}

	// 监听成功
	fmt.Println("启动 TCP-framework 服务 成功", s.Name, "success Listening...")

	// TODO server.go应该有一个自定生成ID的方法
	var cid uint32
	cid = 0

	// 3.阻塞的等待客户端链接业务（读写）
	for true {
		// 3.1 阻塞等待客户端建立链接请求
		tcpConn, err := listenTCP.AcceptTCP()
		if err != nil {
			logrus.Error("Accept err", err)
			continue
		}
		// 3.2 TODO Server.Start() 设置服务器最大链接控制，如果超过最大链接，则关闭此新的链接
		if s.ConnMgr.Len() >= tcp_utils.GlobalObj.MaxConn {
			tcpConn.Close()
			continue
		}
		// 3.3 处理该新链接请求的业务 方法， 此时应该有Handler 和 conn是绑定的
		dealConn := NewConnection(s, tcpConn, cid, s.MessageHandle)
		cid++

		// 3.4 启动当前链接的处理业务
		go dealConn.Start()

	}
}

func (s *Server) Stop() {
	logrus.Info("[ServerSTOP] TCP-framework server", s.Name)
	// TODO 将一些服务器的资源，状态或者一些已经开辟的链接信息进行停止或者清除
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	//TODO 做一些启动服务器之后的业务

	// 阻塞 否则主Go退出，listen的go将会退出
	for true {
		time.Sleep(10 * time.Second)
	}
}

// AddRouter 添加路由
func (s *Server) AddRouter(msgId uint32, router tcp_iface.IRouter) {
	s.MessageHandle.AddRouter(msgId, router)
	fmt.Println("Add Router Success!!")
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() tcp_iface.IConnectionManager {
	return s.ConnMgr
}

//调用连接OnConnStop Hook函数

func (s *Server) SetOnConnStart(hookFunc func(connection tcp_iface.IConnection)) {
	s.OnConnStart = hookFunc
}

//调用连接OnConnStop Hook函数

func (s *Server) SetOnConnStop(hookFunc func(connection tcp_iface.IConnection)) {
	s.OnConnStop = hookFunc
}

//调用连接OnConnStop Hook函数

func (s *Server) CallOnConStart(connection tcp_iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----> CallOnConnStart...")
		s.OnConnStart(connection)
	}
}

//调用连接OnConnStop Hook函数

func (s *Server) CallOnConnStop(connection tcp_iface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----> CallOnConnStop...")
		s.OnConnStop(connection)
	}
}
