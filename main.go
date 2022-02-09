package main

import (
	"github.com/fatih/color"
	"jim_gateway/internal/manager"
	"jim_gateway/pkg"
)

func main() {
	go manager.StartWsServer(pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port)
	go manager.StartTcpServer(pkg.Conf.Tcp.Host, pkg.Conf.Tcp.Port)

	if pkg.Conf.Runtime.Mode == manager.ModelGrpc.String() || pkg.Conf.Runtime.Mode == manager.ModelKafka.String() {
		go manager.RegisterGateway(pkg.Conf.Gateway.Id)
		go manager.ReadMessageFromGrpc()
	}

	if pkg.Conf.Runtime.Mode == manager.ModelKafka.String() {
		go func() {
			err := manager.SubscribeSendMessage()
			if err != nil {
				color.Red("dispatch:cmd:login main err:%s", err.Error())
			}
		}()
	}

	select {}
}
