package main

import (
	"github.com/Siu-Ming/common-go/tcp_iface"
	"github.com/Siu-Ming/common-go/tcp_net"
)

type PingRouter struct {
	tcp_net.BaseRouter
}

func (p *PingRouter) BeforeHandle(request tcp_iface.IRequest) {

}

func (p *PingRouter) Handle(request tcp_iface.IRequest) {
	request.GetConnection().SendBuffMsg(0,[]byte("我是处理消息1的方法"))
}

func (p *PingRouter) AfterHandle(request tcp_iface.IRequest) {

}


type HelloZinxRouter struct {
	tcp_net.BaseRouter
}

func (r *HelloZinxRouter) Handle(request tcp_iface.IRequest) {
	request.GetConnection().SendBuffMsg(1, []byte("ping...ping...ping..."))

}