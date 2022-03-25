package manager

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"jim_gateway/pkg"
	"net"
)

type Tcp struct {
	conn   *net.Conn
	err    error
}

func StartTcpServer(host string, port uint,ctx context.Context) error {
	clientManager := GetClientManagerInstance()
	go clientManager.Loop()
	defer func() {
		clientManager.Close()
	}()
	var address = fmt.Sprintf("%s:%d", host, port)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	listen, err1 := net.ListenTCP("tcp4", addr)
	if err1 != nil {
		return err1
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				color.Red("close tcp")
				return
			}
		}
	}()

	for {
		conn, err2 := listen.AcceptTCP()
		if err2 != nil {
			continue
		}
		client := NewTcpClient(0, "", conn, pkg.Conf.Gateway.Id,pkg.Conf.Runtime.Mode)
		clientManager.Connect(client)
	}

}
