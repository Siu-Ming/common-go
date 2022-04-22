package main

import (
	"fmt"
	"github.com/Siu-Ming/common-go/tcp_iface"
	"github.com/Siu-Ming/common-go/tcp_net"
)

func main() {
	server := tcp_net.NewServer("Ming")

	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)

	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloZinxRouter{})

	server.Serve()
}

// DoConnectionBegin 创建链接前执行
func DoConnectionBegin(conn tcp_iface.IConnection) {
	fmt.Println("Do ConnectionBegin is Called ...")

	conn.SetProperty("Name", "Aceld")
	conn.SetProperty("Home", "https://www.baidu.com")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}

// DoConnectionLost 断开链接后执行
func DoConnectionLost(connection tcp_iface.IConnection) {

	Name, err := connection.GetProperty("Name")
	if err == nil {
		fmt.Println("Conn Name=", Name)
	}
	Home, err := connection.GetProperty("Home")
	if err == nil {
		fmt.Println("Conn Home=", Home)
	}

	fmt.Println("DoConnectionLost is Called...")
}
