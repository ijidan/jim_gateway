package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"jim_message/pkg"
	"net/url"
	"os"
	"sync"
)

var id = flag.Int("id", 0, "client id")
var mutex = sync.Mutex{}

func connect(u url.URL, clientId string, userId int) *websocket.Conn {
	urlStr := fmt.Sprintf("%s?client_id=%s&user_id=%d", u.String(), clientId, userId)
	clientConn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		logrus.Fatalf("connect err:%v", err.Error())
		return nil
	}
	go ReadMessage(clientConn)
	return clientConn
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
			contentStr:=string(messageContent)
			if contentStr!="ping"{
				color.Yellow("received message:%s", string(messageContent))
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
	clientId := fmt.Sprintf("client_%d", *id)
	clientConn:=connect(u, clientId, *id)

	reader:=bufio.NewReader(os.Stdin)
	for  {
		ipt,err:=reader.ReadString(byte('\n'))
		if err!=nil{
			color.Red("read input content error:%s",err.Error())
		}
		err2 := clientConn.WriteMessage(websocket.TextMessage, []byte(ipt))
		if err2!=nil{
			color.Red("write message err:%s",err2.Error())
		}
	}
}
