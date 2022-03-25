package manager

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"jim_gateway/pkg"
	"net/http"
)

//ws类
type webSocket struct {
	upGrader *websocket.Upgrader
	conn     *websocket.Conn
	err      error
}

var upGrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//处理连接
func (ws webSocket) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		ws.err = err
		return
	}
	ws.conn = conn
	clientId := r.FormValue("client_id")
	userId := cast.ToUint64(r.FormValue("user_id"))
	if clientId == "" || userId == 0 {
		err := conn.Close()
		if err != nil {
			return
		}
	}
	client := NewWsClient(userId, clientId,conn,pkg.Conf.Gateway.Id,pkg.Conf.Runtime.Mode)

	clientManager := GetClientManagerInstance()
	clientManager.Connect(client)

	if pkg.Conf.Runtime.Mode==ModelGrpc.String()||pkg.Conf.Runtime.Mode==ModelKafka.String(){
		RegisterClient(pkg.Conf.Gateway.Id,clientId)
	}
	client.Send([]byte("welcome"))
}

func StartWsServer(host string, port uint,ctx context.Context)error {
	clientManager := GetClientManagerInstance()
	go clientManager.Loop()
	defer func() {
		clientManager.Close()
	}()
	var address = fmt.Sprintf("%s:%d", host, port)
	ws := webSocket{
		upGrader: upGrader,
	}
	httpMutex := http.NewServeMux()
	httpMutex.HandleFunc("/",ws.handleConnection)
	httpServer := http.Server{
		Addr:    address,
		Handler: httpMutex,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				color.Red("close websocket")
				_ = httpServer.Shutdown(context.Background())
				return
			}
		}
	}()

	return httpServer.ListenAndServe()





}
