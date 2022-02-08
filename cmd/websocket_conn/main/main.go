package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"jim_message/pkg"
	"net/url"
	"runtime"
	"sync"
	"time"
)

var num = flag.Int("num", 1, "conn num")
var id = flag.Int("id", 0, "client id")
var mutex = sync.Mutex{}

func connect(u url.URL, clientId string, userId int) {
	urlStr := fmt.Sprintf("%s?client_id=%s&user_id=%d", u.String(), clientId, userId)
	clientConn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		logrus.Fatalf("connect err:%v", err.Error())
		return
	}
	go ReadMessage(clientConn)
	go WriteMessage(clientConn, clientId)
	cnt := runtime.NumGoroutine()
	color.Magenta("-----------------------client:【%s】,goroutine cnt:【%d】------------------", clientId, cnt)
}

func ReadMessage(clientConn *websocket.Conn) {
	for {
		_, messageContent, err1 := clientConn.ReadMessage()
		if err1 != nil {
			if err1 == io.EOF {
				continue
			} else {
				Close(clientConn)
				color.Red("connect err:%v", err1.Error())
				return
			}
		} else {
			color.Yellow("received message:%s", string(messageContent))
		}
	}
}

func WriteMessage(clientConn *websocket.Conn, clientId string) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			message := fmt.Sprintf("from %s@%d", clientId, time.Now().Unix())
			err2 := clientConn.WriteMessage(websocket.TextMessage, []byte(message))

			if err2 != nil {
				ticker.Stop()
				Close(clientConn)
				color.Red("send message error:%s", err2.Error())
				return
			} else {
				color.Green("send message success:%s", message)
			}
		}
	}
}
func Close(clientConn *websocket.Conn) {
	mutex.Lock()
	_ = clientConn.Close()
	mutex.Unlock()
}
func main() {
	flag.Parse()
	var address = fmt.Sprintf("%s:%d", pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port)
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	if *id > 0 {
		clientId := fmt.Sprintf("client_%d", *id)
		go connect(u, clientId, *id)
	} else {
		for i := 0; i < *num; i++ {
			clientId := fmt.Sprintf("client_id_%d", i+1)
			go connect(u, clientId, i+1)
		}
	}
	logrus.Println("waiting...")
	select {}
}
