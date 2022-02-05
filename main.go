package main

import (
	"jim_message/internal/manager"
	"jim_message/pkg"
)

func main() {
	go manager.StartWsServer(pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port)
	go manager.StartTcpServer(pkg.Conf.Tcp.Host, pkg.Conf.Tcp.Port)
	select {}
}
