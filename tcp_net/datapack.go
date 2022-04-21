package tcp_net

import (
	"TCP-framework_V1.0/tcp_iface"
	"TCP-framework_V1.0/tcp_utils"
	"bytes"
	"encoding/binary"
	"errors"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

// Pack 封包
func (dp DataPack) Pack(msg tcp_iface.IMessage) ([]byte, error) {
	//写msgID
	buffer := bytes.NewBuffer([]byte{})

	//写msgID
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// UnPack 拆包
func (dp DataPack) UnPack(binaryData []byte) (tcp_iface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	readerData := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(readerData, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(readerData, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if tcp_utils.GlobalObj.MaxPacketSize > 0 && msg.DataLen > tcp_utils.GlobalObj.MaxPacketSize {
		return nil, errors.New("Too large msg data recieved")
	}

	return msg, nil
}
