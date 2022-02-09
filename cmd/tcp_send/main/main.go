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
)

var id=flag.Int("id",0,"client id")
var mutex = sync.Mutex{}

func connect(address *net.TCPAddr, clientId string, userId int) *net.TCPConn{
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
			message:=string(data[manager.HeaderFlagLen+manager.HeaderBodyLen:])
			if message!="ping"{
				color.Yellow("received message:%s",message)
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
	clientId:=cast.ToString(*id)
	clientConn:=connect(addr, clientId, *id)
	reader:=bufio.NewReader(os.Stdin)
	for  {
		ipt,err:=reader.ReadString(byte('\n'))
		if err!=nil{
			color.Red("read input content error:%s",err.Error())
		}

		content, _ := manager.Pack(ipt)
		_, err2 := clientConn.Write(content)
		if err2!=nil{
			color.Red("write message err:%s",err2.Error())
		}
	}
}
