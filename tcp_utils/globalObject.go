package tcp_utils

import (
	"encoding/json"
	"github.com/Siu-Ming/common-go/tcp_iface"
	"io/ioutil"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 zinx.json来配置
*/

type GlobalObject struct {
	//当前服务器主机允许的最大链接个数
	TcpServer tcp_iface.IServer
	//当前服务器主机IP
	Host string
	//当前服务器主机IP
	Port int
	//当前服务器主机IP
	Name string
	//当前版本号
	Version string

	//需数据包的最大值
	MaxPacketSize uint32
	//当前服务器主机允许的最大链接个数
	MaxConn int

	WorkerPoolSize uint32

	MaxWorkerTaskLen uint32

	ConfFilePath string

	MaxMsgChanLen int
}

var GlobalObj *GlobalObject

func (g *GlobalObject) ReloadFile() {
	readFile, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(readFile, &GlobalObj)
	if err != nil {
		panic(err)
	}
}

/*
	提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalObj = &GlobalObject{
		Name:             "TCP-framework",
		Version:          "V1.0",
		Port:             8999,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/config.json",
		MaxMsgChanLen:    10,
	}

	// 从配置文件中加载一些用户配置的参数
	// GlobalObj.ReloadFile()
}
