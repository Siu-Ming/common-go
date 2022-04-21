package tcp_net

import "TCP-framework_V1.0/tcp_iface"

type Request struct {
	// 已经和客户端建立好的 链接
	conn tcp_iface.IConnection
	// 客户端请求的数据
	// data []byte
	// 客户端请求数据
	msg tcp_iface.IMessage
}

// GetConnection 获取请求连接信息
func (r *Request) GetConnection() tcp_iface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgId GetMsgID 获取请求的消息的ID
func (r Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
