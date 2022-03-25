package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"jim_gateway/config"
	"jim_gateway/internal/manager"
	"jim_gateway/pkg"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func buildTable(config config.Config) *tablewriter.Table {
	tcpAddress := fmt.Sprintf("%s:%d", config.Tcp.Host, config.Tcp.Port)
	wsAddress := fmt.Sprintf("%s:%d", config.Websocket.Host, config.Websocket.Port)
	grpcAddress := fmt.Sprintf("%s:%d", config.Rpc.Host, config.Rpc.Port)
	pprofAddress := fmt.Sprintf("%s:%d", config.Pprof.Host, config.Pprof.Port)


	data := [][]string{
		[]string{"1", "Application", "Jim_Gateway"},
		[]string{"2", "Tcp", tcpAddress},
		[]string{"3", "Grpc", grpcAddress},
		[]string{"4", "Websocket", wsAddress},
		[]string{"5", "Pprof", pprofAddress},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Info", "Desc"})
	table.AppendBulk(data)
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	return table
}

func main() {
	defer pkg.Close()
	table:=buildTable(*pkg.Conf)
	table.Render()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err:=manager.StartWsServer(pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port,ctx)
		if err!=nil{
			cancel()
			pkg.Logger.Fatalf("run websocket server:%s", err.Error())
		}
	}()
	go func() {
		err:= manager.StartTcpServer(pkg.Conf.Tcp.Host, pkg.Conf.Tcp.Port,ctx)
		if err!=nil{
			cancel()
			pkg.Logger.Fatalf("run tcp server:%s", err.Error())
		}
	}()
	go func() {
		err := manager.RunPprof(*pkg.Conf, ctx)
		if err != nil {
			cancel()
			pkg.Logger.Fatalf("run pprof server:%s", err.Error())
		}
	}()

	if pkg.Conf.Runtime.Mode == manager.ModelGrpc.String() || pkg.Conf.Runtime.Mode == manager.ModelKafka.String() {
		go manager.RegisterGateway(pkg.Conf.Gateway.Id)
		go manager.ReadMessageFromGrpc()
	}



	//if pkg.Conf.Runtime.Mode == manager.ModelKafka.String() {
	//	go func() {
	//		err := manager.SubscribeSendMessage()
	//		if err != nil {
	//			color.Red("dispatch:cmd:login main err:%s", err.Error())
	//		}
	//	}()
	//}
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGILL, syscall.SIGQUIT, syscall.SIGTERM)
	<-ch
	color.Red("closing ...")
	cancel()
	color.Red("closed")
}
