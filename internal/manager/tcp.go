package manager

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"jim_gateway/pkg"
	"net"
)

type Tcp struct {
	conn   *net.Conn
	err    error
}

func StartTcpServer(host string, port uint) {
	clientManager := GetClientManagerInstance()
	go clientManager.Loop()
	defer func() {
		clientManager.Close()
	}()
	var address = fmt.Sprintf("%s:%d", host, port)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		logrus.Fatalf("tcp resolve error:%s", err.Error())
	}
	listen, err1 := net.ListenTCP("tcp4", addr)
	if err1 != nil {
		logrus.Fatalf("tcp listen error:%s", err1.Error())
	}
	for {
		conn, err2 := listen.AcceptTCP()
		if err2 != nil {
			logrus.Fatalf("tcp accept error:%s", err2.Error())
		}
		client := NewTcpClient(0, "", conn, pkg.Conf.Gateway.Id,pkg.Conf.Runtime.Mode)
		clientManager.Connect(client)
	}

}
