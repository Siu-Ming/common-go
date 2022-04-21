package tcp_iface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(message IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
