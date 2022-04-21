package tcp_iface

/**
链接管理抽象层
*/

type IConnectionManager interface {
	// Add 添加链接
	Add(conn IConnection)
	// Remove 删除链接
	Remove(conn IConnection)
	// Get 利用ConnID获取链接
	Get(connId uint32) (IConnection, error)
	// Len 获取当前链接长度
	Len() int
	// ClearConn 删除并停止所有链接
	ClearConn()
}
