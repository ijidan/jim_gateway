package manager

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"jim_gateway/pkg"
	"net"
)

type Tcp struct {
	server *Server
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
		var server = NewServer(address)
		client := NewTcpClient(0, "", server, conn, pkg.Conf.Runtime.Mode)
		server.AddClient(client)
		clientManager.Connect(client)
	}

}
