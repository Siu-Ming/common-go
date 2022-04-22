package main

import (
	"fmt"
	"github.com/Siu-Ming/common-go/tcp_net"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		return
	}

	dataPack := tcp_net.NewDataPack()

	msg1 := tcp_net.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	msg2 := tcp_net.Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'m', 'i', 'n', 'g'},
	}

	sendData1, _ := dataPack.Pack(&msg1)

	sendData2, _ := dataPack.Pack(&msg2)

	data := append(sendData1, sendData2...)

	fmt.Println(string(data))
	conn.Write(data)

	select {

	}

}
