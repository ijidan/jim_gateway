package main

import (
	"github.com/fatih/color"
	"jim_message/internal/dispatch"
	"jim_message/internal/manager"
	"jim_message/pkg"
)

func main() {
	go manager.StartWsServer(pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port)
	go manager.StartTcpServer(pkg.Conf.Tcp.Host, pkg.Conf.Tcp.Port)
	go func() {
		err := dispatch.SubscribeCmdLogin()
		if err != nil {
			color.Red("dispatch:cmd:login main err:%s",err.Error())
		}
	}()
	select {}
}
