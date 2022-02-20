package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"io"
	"jim_gateway/internal/manager"
	"jim_gateway/pkg"
	"net"
	"os"
	"sync"
	"time"
)

var id = flag.Int("id", 0, "client id")
var mutex = sync.Mutex{}

func connect(address *net.TCPAddr, clientId string, userId int) *net.TCPConn {
	clientConn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		logrus.Fatalf("connect err:%v", err.Error())
		return nil
	}
	go ReadMessage(clientConn)
	return clientConn
}

func ReadMessage(clientConn *net.TCPConn) {
	for {
		reader := bufio.NewReader(clientConn)
		header, err1 := reader.Peek(manager.BusinessHeaderFlagLen + manager.BusinessHeaderCmdLen + manager.BusinessHeaderRequestIdLen + manager.BusinessHeaderContentLen)
		if err1 != nil {
			if err1 == io.EOF {
				continue
			} else {
				Close(clientConn)
				return
			}
		}
		color.Yellow("1111111111111")
		if !bytes.HasPrefix(header, []byte(manager.BusinessHeaderFlag)) {
			Close(clientConn)
			return
		}
		color.Yellow("22222222222222")
		headerBody := header[manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen:]
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
		color.Yellow("333333333333333333")
		if int32(reader.Buffered()) < manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen+manager.BusinessHeaderContentLen+bodyLen {
			color.Yellow("4444444444444444")
			continue
		}
		data := make([]byte, manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen+manager.BusinessHeaderContentLen+bodyLen)
		_, err3 := reader.Read(data)
		color.Yellow("555555555:%s", string(data))
		if err3 != nil {
			if err3 == io.EOF {
				continue
			} else {
				Close(clientConn)
				return
			}
		} else {
			requestId := data[manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen : manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen]
			color.Yellow("received requestId:%d", manager.BytesToUint32(requestId))

			message := string(data[manager.BusinessHeaderFlagLen+manager.BusinessHeaderCmdLen+manager.BusinessHeaderRequestIdLen+manager.BusinessHeaderContentLen:])
			if message != "ping" {
				color.Yellow("received message:%s", message)
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
	clientId := cast.ToString(*id)
	clientConn := connect(addr, clientId, *id)
	reader := bufio.NewReader(os.Stdin)
	for {
		ipt, err := reader.ReadString(byte('\n'))
		if err != nil {
			color.Red("read input content error:%s", err.Error())
		}
		requestId := uint32(time.Now().Second())
		content, _ := manager.BusinessPack(manager.BusinessCmdC2C, requestId, ipt)
		_, err2 := clientConn.Write(content)
		if err2 != nil {
			color.Red("write message err:%s", err2.Error())
		}
	}
}
