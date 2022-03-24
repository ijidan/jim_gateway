package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"jim_gateway/internal/manager"
	"jim_gateway/pkg"
	"net"
	"strconv"
	"sync"
	"time"
)

var num = flag.Int("num", 1, "conn num")
var mutex = sync.Mutex{}

func connect(address *net.TCPAddr, clientId string, token string) *net.TCPConn {
	clientConn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		logrus.Fatalf("connect err:%v", err.Error())
		return nil
	}
	Login(clientConn,clientId,token)
	go ReadMessage(clientConn)
	go HeartBeat(clientConn)
	return clientConn
}

func Login(clientConn *net.TCPConn,clientId string, token string)  {
	authMessage:=manager.AuthLoginMessage{
		ClientId: clientId,
		Token:    token,
	}
	messageContent,_:=json.Marshal(authMessage)
	content, _ := manager.BusinessPack(manager.BusinessCmdAuthLogin, 0, string(messageContent))
	_, err := clientConn.Write(content)
	if err!=nil{
		Close(clientConn,err)
	}
}
func ReadMessage(clientConn *net.TCPConn) {
	err99 := clientConn.SetDeadline(time.Now().Add(1 * time.Minute))
	if err99 != nil {
		Close(clientConn,err99)
		return
	}
	for {
		reader := bufio.NewReader(clientConn)
		header, err1 := reader.Peek(manager.BusinessHeaderFlagLen + manager.BusinessHeaderCmdLen + manager.BusinessHeaderRequestIdLen + manager.BusinessHeaderContentLen)
		if err1 != nil {
			if err1 == io.EOF {
				continue
			} else {
				Close(clientConn,err1)
				return
			}
		}

		if !bytes.HasPrefix(header, []byte(manager.BusinessHeaderFlag)) {
			color.Red("receiver header:%s,%s",header,[]byte(manager.BusinessHeaderFlag))
			Close(clientConn,errors.New("header flag error"))
			return
		}
		headerBody := header[manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen:]
		buffer := bytes.NewBuffer(headerBody)
		var bodyLen int32
		err2 := binary.Read(buffer, binary.BigEndian, &bodyLen)
		if err2 != nil {
			if err2 == io.EOF {
				continue
			} else {
				Close(clientConn,err2)
				return
			}
		}
		if int32(reader.Buffered()) < manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen+manager.BusinessHeaderContentLen+bodyLen {
			continue
		}
		data := make([]byte, manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen+manager.BusinessHeaderContentLen+bodyLen)
		_, err3 := reader.Read(data)
		if err3 != nil {
			if err3 == io.EOF {
				continue
			} else {
				Close(clientConn,err3)
				return
			}
		} else {
			headerCmd := data[manager.BusinessHeaderFlagLen : manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen]
			if bytes.Compare(headerCmd, []byte(manager.BusinessCmdPong)) == 0 {
				err88 := clientConn.SetDeadline(time.Now().Add(1 * time.Minute))
				if err88 != nil {
					Close(clientConn,errors.New("header cmd error"))
					return
				}
			}
			//requestId := data[manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen : manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen]
			//color.Yellow("received requestId:%d", manager.BytesToUint32(requestId))
			message := string(data[manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen+manager.BusinessHeaderContentLen:])
			if message != "pong" {
				color.Yellow("received message:%s", string(message))
			}
		}
	}
}

func HeartBeat(clientConn *net.TCPConn) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			content, _ := manager.BusinessPack(manager.BusinessCmdPing, 0, "ping")
			//color.Green("heart beat:ping")
			_, err := clientConn.Write(content)
			if err != nil {
				color.Red("heart beat error:%s", err.Error())
			}
		}
	}
}


func Close(clientConn *net.TCPConn,err error) {
	mutex.Lock()
	color.Red("client close:%s",err.Error())
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
	for i := 0; i < *num; i++ {
		clientId := fmt.Sprintf("%d", i+1)
		clientId=strconv.Itoa(i+1)
		var clientConn *net.TCPConn
		go func(conn *net.TCPConn) {
			conn=connect(addr, clientId, clientId)
		}(clientConn)
		defer func(conn *net.TCPConn) {
			_=conn.Close()
		}(clientConn)
	}
	logrus.Println("waiting...")
	select {}

}
