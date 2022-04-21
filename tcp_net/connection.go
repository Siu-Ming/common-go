package tcp_net

import (
	"TCP-framework_V1.0/tcp_iface"
	"TCP-framework_V1.0/tcp_utils"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// 当前Conn属于哪一个Server
	TcpServer tcp_iface.IServer

	// 当前链接的socketTCP套接字
	Conn *net.TCPConn
	// 当前链接的ID 也可以称为SessionID，ID全局唯一
	ConnID uint32
	// 当前链接的关闭状态
	isClosed bool

	//该连接的处理方法router
	// Router tcp_iface.IRouter

	// 消息管理 和对应的业务API关系
	MessageHandle tcp_iface.IMessageHandle

	// 告知该链接已经退出/停止的channel
	ExitBuffChannel chan bool

	// 无缓冲管道， 用于读，写两个Goroutine之间通信
	msgChan chan []byte

	// 有缓冲管道 用于读，写两个Goroutine之间消息通信
	msgBuffChannel chan []byte

	// 链接属性
	property map[string]interface{}
	// 保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewConnection(server tcp_iface.IServer, conn *net.TCPConn, connId uint32, handle tcp_iface.IMessageHandle) *Connection {
	connection := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connId,
		isClosed:  false,
		//Router:          router,
		ExitBuffChannel: make(chan bool, 1),
		MessageHandle:   handle,
		msgChan:         make(chan []byte),
		msgBuffChannel:  make(chan []byte, tcp_utils.GlobalObj.MaxMsgChanLen),
		property:        make(map[string]interface{}),
	}
	// 将新建的Conn添加到链接管理中
	connection.TcpServer.GetConnMgr().Add(connection)
	return connection
}

// StartWriter 写消息Goroutine，将用户数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")

	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	// 阻塞等待channel消息  进行回写给客户端
	for true {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				logrus.Error("Send Data error:", err, "Conn Writer exit!")
				return
			}
		case data, ok := <-c.msgBuffChannel:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:,", err, "Conn Writer exit")
					return
				}
			} else {
				break
				fmt.Println("msgBuffChan is Closed")
			}
		case <-c.ExitBuffChannel:
			// Conn已经关闭
			// 代表牌Reader已经退出，此时Writer也要退出
			return
		}
	}

}

func (c *Connection) StartReader() {
	fmt.Println("[reader Goroutine is running]")
	defer logrus.Info(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()

	for true {

		// 创建拆包解包对象
		dp := NewDataPack()

		// 读取客户端Msg head
		headData := make([]byte, dp.GetHeadLen())

		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			logrus.Error("read msg head error", err)
			c.ExitBuffChannel <- true
			break
		}

		// 拆包, 得到msgid 和 datalen 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			logrus.Error("unpack error", err)
			c.ExitBuffChannel <- true
			break
		}

		// 根据dataLen读取data 放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				logrus.Error("read msg data error", err)
				c.ExitBuffChannel <- true
				break
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的Request
		req := Request{
			conn: c,
			msg:  msg,
		}
		if tcp_utils.GlobalObj.WorkerPoolSize > 0 {
			// 已启动工作池机制，将消息交给worker处理
			c.MessageHandle.SendMessageToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID 找到对应处理api的业务
			go c.MessageHandle.DoMsgHandler(&req)
		}

	}
}

func (c *Connection) Start() {
	// 开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	// 启动当前链接写数据的业务
	go c.StartWriter()

	// 按照用户传递进来的创建链接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConStart(c)

	for true {
		select {
		case <-c.ExitBuffChannel:
			// 得到退出消息，不再阻塞
			return
		}
	}
}

func (c *Connection) Stop() {
	logrus.Info("Conn Stop()...ConnID=", c.ConnID)
	// 1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	// true为关闭状态
	c.isClosed = true

	// 如果用户注册了改该链接的关闭回调业务，呢么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()
	// 通知从缓冲队列读数据业务，该链接已经关闭
	c.ExitBuffChannel <- true
	// 将链接从管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)
	// 关闭该连接全部管道
	close(c.ExitBuffChannel)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection close when send msg")
	}
	// 将data封装 并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		logrus.Error("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	// 消息发送给管道
	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		logrus.Error("Pack error msg id=", msgId)
		return err
	}
	// 写回客户端
	c.msgBuffChannel <- msg
	return nil
}

// SetProperty 设置属性链接
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取属性链接
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除属性链接
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
