package manager

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"net/http"
)

//ws类
type webSocket struct {
	upGrader *websocket.Upgrader
	server   *Server
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
	client := NewWsClient(userId, clientId, ws.server, conn)
	ws.server.AddClient(client)

	//clientManager := GetClientManagerInstance()
	//clientManager.Connect(client)

	message := NewMessage(0, userId, []byte("welcome"))
	client.Send(message)

}

func StartWsServer(host string, port uint) {
	clientManager := GetClientManagerInstance()
	go clientManager.Loop()
	defer func() {
		clientManager.Close()
	}()
	var address = fmt.Sprintf("%s:%d", host, port)
	var server = NewServer(address)
	ws := webSocket{
		upGrader: upGrader,
		server:   server,
	}
	http.HandleFunc("/", ws.handleConnection)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		logrus.Fatalf("websocket listen error:%s", err.Error())
	}

}
