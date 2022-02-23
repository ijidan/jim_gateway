package manager

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"jim_gateway/internal/jim_proto/proto_build"
	"net"
	"sync"
	"time"
)

type ConnType uint8

const (
	ConnTypeTcp ConnType = iota
	ConnTypeWs
)

type Client struct {
	clientId  string //client id
	userId    uint64 //user id
	connType  ConnType
	tcpConn   *net.TCPConn    //tcp conn
	wsConn    *websocket.Conn //websocket conn
	readCh    chan []byte     //read message channel
	writeCh   chan []byte     //send message channel
	closeCh   chan byte
	mutex     sync.Mutex
	isRunning bool
	gatewayId string
	mode      string
}

func (c *Client) GetClientId() string {
	return c.clientId
}

func (c *Client) ReadMessage() {
	clientManager := GetClientManagerInstance()
	if c.isRunning == false {
		return
	}
	if c.connType == ConnTypeTcp {
		err99 := c.tcpConn.SetDeadline(time.Now().Add(1 * time.Minute))
		if err99 != nil {
			c.Close(err99)
			return
		}
	} else {
		err99 := c.wsConn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		if err99 != nil {
			c.Close(err99)
			return
		}
	}
	for {
		if c.isRunning == false {
			return
		}
		var messageContent []byte
		//var err error
		//if c.connType == ConnTypeTcp {

		reader := bufio.NewReader(c.tcpConn)
		header, err0 := reader.Peek(BusinessHeaderFlagLen + BusinessHeaderCmdLen + BusinessHeaderRequestIdLen + BusinessHeaderContentLen)
		if err0 != nil {
			if err0 == io.EOF {
				continue
			} else {
				c.Close(err0)
				return
			}
		}

		if !bytes.HasPrefix(header, []byte(BusinessHeaderFlag)) {
			c.Close(errors.New("header flag error"))
			return
		}

		headerBody := header[BusinessHeaderFlagLen+BusinessHeaderCmdLen+BusinessHeaderRequestIdLen:]
		buffer := bytes.NewBuffer(headerBody)
		var bodyLen int32
		err2 := binary.Read(buffer, binary.BigEndian, &bodyLen)
		if err2 != nil {
			if err2 == io.EOF {
				continue
			} else {
				c.Close(err2)
				return
			}
		}

		if int32(reader.Buffered()) < BusinessHeaderFlagLen+BusinessHeaderCmdLen+BusinessHeaderRequestIdLen+BusinessHeaderContentLen+bodyLen {
			continue
		}
		data := make([]byte, BusinessHeaderFlagLen+BusinessHeaderCmdLen+BusinessHeaderRequestIdLen+BusinessHeaderContentLen+bodyLen)
		_, err3 := reader.Read(data)
		if err3 != nil {
			if err3 == io.EOF {
				continue
			} else {
				c.Close(err3)
				return
			}
		}
		headerCmd := string(data[BusinessHeaderFlagLen : BusinessHeaderFlagLen+BusinessHeaderCmdLen])
		requestId := BytesToUint32(data[BusinessHeaderFlagLen+BusinessHeaderCmdLen : BusinessHeaderFlagLen+BusinessHeaderCmdLen+BusinessHeaderRequestIdLen])
		messageContent = data[BusinessHeaderFlagLen+BusinessHeaderCmdLen+BusinessHeaderRequestIdLen+BusinessHeaderContentLen:]

		switch headerCmd {
		case BusinessCmdPing:
			content, _ := BusinessPack(BusinessCmdPong, 0, "pong")
			_, err4 := c.tcpConn.Write(content)
			if err4 != nil {
				c.Close(err4)
				return
			}
		case BusinessCmdAuthLogin:
			var authLogin AuthLoginMessage
			_=json.Unmarshal(messageContent,&authLogin)
			clientId:=authLogin.ClientId
			token:=authLogin.Token
			c.clientId=clientId
			c.userId=ParseToken(token)
			_, errLogin := RegisterClient(c.gatewayId, c.clientId)
			if errLogin != nil {
				c.Close(errLogin)
				return
			}
			clientManager.Connect(c)
			content, _ := BusinessPack(BusinessCmdAuthSuccess, 0, "auth success")
			_, err4 := c.tcpConn.Write(content)
			if err4 != nil {
				c.Close(err4)
				return
			}
		case BusinessCmdAuthLogout:
			_, errLogout := UnRegisterClient(c.gatewayId, c.clientId)
			if errLogout != nil {
				c.Close(errLogout)
				return
			}
			content, _ := BusinessPack(BusinessCmdServerClose, 0, "client close")
			_, err5 := c.tcpConn.Write(content)
			if err5 != nil {
				c.Close(err5)
				return
			}
		default:
			color.Yellow("received message:%s",string(messageContent))
			req := &proto_build.SendMessageRequest{
				GatewayId: c.gatewayId,
				Cmd:       headerCmd,
				RequestId: requestId,
				Data:      messageContent,
			}

			color.Red("grpc send gateway...............%s", c.gatewayId)
			color.Red("grpc send cmd...............%s", headerCmd)
			color.Red("grpc send requestId...............%d",requestId)
			color.Red("grpc send data...............%s", string(messageContent))

			sendClient := GetGatewayServiceSendMessageClient()
			errSend1 := sendClient.Send(req)
			if errSend1 != nil {
				color.Red("send client send message error:%s", errSend1.Error())
				c.Close(errSend1)
				return
			}
		}

		//}
		//else {
		//	_, messageContent, err = c.wsConn.ReadMessage()
		//	if err != nil {
		//		if err == io.EOF {
		//			continue
		//		} else {
		//			color.Red("message read error:" + err.Error())
		//			c.Close(err)
		//			return
		//		}
		//	}
		//}
		err88 := c.tcpConn.SetDeadline(time.Now().Add(1 * time.Minute))
		if err88 != nil {
			c.Close(err88)
			return
		}
		if string(messageContent)!="ping"{
			color.Yellow("message received:%s", string(messageContent))
		}

		//if c.mode == ModeLocal.String() {
		//	var data json.RawMessage
		//	msg := ClientMessage{
		//		Data: &data,
		//	}
		//	if err3 := json.Unmarshal(messageContent, &msg); err3 != nil {
		//		color.Red("parse message err:%s", err3.Error())
		//	}
		//	switch msg.Cmd {
		//	case "auth.req":
		//		clientId, userId := ParseAuthReqMessage(data)
		//		c.clientId = clientId
		//		c.userId = userId
		//		clientManager.Connect(c)
		//	case "chat.c2c.txt":
		//		ParseC2CTxtMessage(data, messageContent)
		//	}
		//}
	}
}
func (c *Client) WriteMessage() {
	//ticker := time.NewTicker(1 * time.Second)
	for {
		if c.isRunning == false {
			return
		}
		select {
		case <-c.closeCh:
			close(c.readCh)
			close(c.writeCh)
			//ticker.Stop()
			return
		//case <-ticker.C:
		//	message := NewMessage(0, c.userId, []byte("ping"))
		//	var err error
		//	if c.connType == ConnTypeTcp {
		//
		//		requestId := uint32(time.Now().Unix())
		//		content, _ := BusinessPack(BusinessCmdPing, requestId, "PING")
		//		color.Green("message send:%d", requestId)
		//		_, err = c.tcpConn.Write(content)
		//	} else {
		//		err = c.wsConn.WriteMessage(websocket.TextMessage, message.data)
		//	}
		//	if err != nil {
		//		c.Close(err)
		//		color.Red("send message error:%s", err.Error())
		//		return
		//	}
		case message, ok := <-c.writeCh:
			if ok {
				var err error
				if c.connType == ConnTypeTcp {
					color.Cyan("send message:%s,%s",string(message),c.clientId)
					_, err = c.tcpConn.Write(message)
				} else {
					err = c.wsConn.WriteMessage(websocket.TextMessage, message)
				}
				if err != nil {
					c.Close(err)
					logrus.Println("send message error:%s", err.Error())
					return
				}
			} else {
				c.Close(errors.New("write ch error"))
				return
			}
		}
	}

}

func (c *Client) Send(message []byte) {
	if c.isRunning {
		color.Green("message send:%s", string(message))
		c.writeCh <- message
	}
}

func (c *Client) Close(err error) {
	c.mutex.Lock()
	if c.isRunning {
		if c.connType == ConnTypeTcp {
			_ = c.tcpConn.Close()
		} else {
			_ = c.wsConn.Close()
		}
		c.isRunning = false
		clientManager := GetClientManagerInstance()
		clientManager.DisConnect(c)

	}
	c.mutex.Unlock()
	logrus.Println("close triggered:%s", err.Error())
}

func NewWsClient(userId uint64, clientId string, conn *websocket.Conn, gatewayId string, mode string) *Client {
	client := &Client{
		clientId:  clientId,
		userId:    userId,
		connType:  ConnTypeWs,
		tcpConn:   nil,
		wsConn:    conn,
		readCh:    make(chan []byte, 1000),
		writeCh:   make(chan []byte, 1000),
		closeCh:   make(chan byte),
		mutex:     sync.Mutex{},
		isRunning: true,
		gatewayId: gatewayId,
		mode:      mode,
	}
	go client.ReadMessage()
	go client.WriteMessage()
	return client
}

func NewTcpClient(userId uint64, clientId string, conn *net.TCPConn, gatewayId string, mode string) *Client {
	client := &Client{
		clientId:  clientId,
		userId:    userId,
		connType:  ConnTypeTcp,
		tcpConn:   conn,
		wsConn:    nil,
		readCh:    make(chan []byte, 1000),
		writeCh:   make(chan []byte, 1000),
		closeCh:   make(chan byte),
		mutex:     sync.Mutex{},
		isRunning: true,
		gatewayId: gatewayId,
		mode:      mode,
	}
	go client.ReadMessage()
	go client.WriteMessage()
	return client
}
