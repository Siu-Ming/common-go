package tcp_net

import (
	"errors"
	"fmt"
	"github.com/Siu-Ming/common-go/tcp_iface"
	"github.com/sirupsen/logrus"
	"sync"
)

type ConnectionManager struct {
	// 管理链接信息
	connections map[uint32]tcp_iface.IConnection
	// 读写锁
	connLock sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]tcp_iface.IConnection),
	}
}

// Add 添加链接
func (connMgr *ConnectionManager) Add(conn tcp_iface.IConnection) {
	// 保护共享资源map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn

	logrus.Info("connection add to ConnManager successfully:conn num=", connMgr.Len())
}

// Remove 删除
func (connMgr *ConnectionManager) Remove(conn tcp_iface.IConnection) {
	// 写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections, conn.GetConnID())

	logrus.Info("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnectionManager) Get(connId uint32) (tcp_iface.IConnection, error) {
	// 读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RLock()
	if conn, ok := connMgr.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnectionManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnectionManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All Connections successfully：conn num=", connMgr.Len())
}
