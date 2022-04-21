package tcp_iface

/**
请求消息封装，抽象接口层
*/

type IMessage interface {
	// GetDataLen 获取消息数据长度
	GetDataLen() uint32
	// GetMsgId 获取消息ID
	GetMsgId() uint32
	// GetData 获取消息内容
	GetData() []byte

	// SetDataLen 设置消息长度
	SetDataLen(uint32)
	// SetMsgId 设置消息ID
	SetMsgId(uint32)
	// SetData 设置消息内容
	SetData([]byte)
}
