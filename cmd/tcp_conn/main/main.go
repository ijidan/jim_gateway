package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"jim_message/internal/manager"
	"jim_message/pkg"
	"net"
	"runtime"
	"sync"
	"time"
)

var num = flag.Int("num", 1, "conn num")
var id=flag.Int("id",0,"client id")
var mutex = sync.Mutex{}

func connect(address *net.TCPAddr, clientId string, userId int) {
	clientConn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		logrus.Fatalf("connect err:%v", err.Error())
		return
	}
	go ReadMessage(clientConn)
	go WriteMessage(clientConn, clientId)
	cnt := runtime.NumGoroutine()
	color.Magenta("-----------------------client:【%s】,goroutine cnt:【%d】------------------", clientId, cnt)
}

func ReadMessage(clientConn *net.TCPConn) {
	for {
		reader := bufio.NewReader(clientConn)
		header, err1 := reader.Peek(manager.HeaderFlagLen + manager.HeaderBodyLen)
		if err1 != nil {
			if err1 == io.EOF {
				continue
			} else {
				Close(clientConn)
				return
			}
		}
		if !bytes.HasPrefix(header, manager.HeaderFlag) {
			Close(clientConn)
			return
		}
		headerBody := bytes.TrimPrefix(header, manager.HeaderFlag)
		buffer := bytes.NewBuffer(headerBody)
		var bodyLen int32
		err2 := binary.Read(buffer, binary.BigEndian, &bodyLen)
		if err2 != nil {
			if err2 == io.EOF {
				continue
			} else {
				Close(clientConn)
				return
			}
		}
		if int32(reader.Buffered()) < manager.HeaderFlagLen+manager.HeaderBodyLen+bodyLen {
			continue
		}
		data := make([]byte, manager.HeaderFlagLen+manager.HeaderBodyLen+bodyLen)
		_, err3 := reader.Read(data)
		if err3 != nil {
			if err3 == io.EOF {
				continue
			} else {
				Close(clientConn)
				return
			}
		} else {
			//color.Yellow("received message:%s", string(data[manager.HeaderFlagLen+manager.HeaderBodyLen:]))
		}
	}
}

func WriteMessage(clientConn *net.TCPConn, clientId string) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			message := fmt.Sprintf("from %s@%d", clientId, time.Now().Unix())
			content, _ := manager.Pack(message)
			_, err2 := clientConn.Write(content)
			if err2 != nil {
				ticker.Stop()
				Close(clientConn)
				//color.Red("send message error:%s", err2.Error())
				return
			} else {
				//color.Green("send message success:%s", message)
			}
		}
	}
}
func Close(clientConn *net.TCPConn) {
	mutex.Lock()
	_ = clientConn.Close()
	mutex.Unlock()
}
func main() {
	flag.Parse()
	var address = fmt.Sprintf("%s:%d", pkg.Conf.Tcp.Host, pkg.Conf.Tcp.Port)
	var addr, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		logrus.Fatalf("tcp resolve err:%s", err.Error())
	}
	if *id>0{
		clientId := fmt.Sprintf("%d", id)
		go connect(addr, clientId, *id)
	}else {
		for i := 0; i < *num; i++ {
			clientId := fmt.Sprintf("client_id_%d", i+1)
			go connect(addr, clientId, i+1)
		}
	}

	logrus.Println("waiting...")
	select {}

}
