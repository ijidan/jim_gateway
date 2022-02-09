package main

import (
	"github.com/fatih/color"
	"jim_gateway/internal/call"
	"jim_gateway/internal/dispatch"
	"jim_gateway/internal/manager"
	"jim_gateway/pkg"
)

func main() {
	go manager.StartWsServer(pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port)
	go manager.StartTcpServer(pkg.Conf.Tcp.Host, pkg.Conf.Tcp.Port)

	if pkg.Conf.Runtime.Mode==manager.ModelGrpc.String(){
		go call.RegisterGateway(pkg.Conf.Gateway.Id)
		go call.ReadMessageFromGrpc()
	}

	if pkg.Conf.Runtime.Mode==manager.ModelKafka.String(){
		go func() {
			err := dispatch.SubscribeCmdLogin()
			if err != nil {
				color.Red("dispatch:cmd:login main err:%s",err.Error())
			}
		}()
	}

	select {}
}
