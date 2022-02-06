package manager

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
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
	clientId  string  //client id
	server    *Server //server
	userId    uint64  //user id
	connType  ConnType
	tcpConn   *net.TCPConn    //tcp conn
	wsConn    *websocket.Conn //websocket conn
	readCh    chan *Message   //read message channel
	writeCh   chan *Message   //send message channel
	closeCh   chan byte
	mutex     sync.Mutex
	isRunning bool
}

func (c *Client) GetClientId() string {
	return c.clientId
}

func (c *Client) GetServer() *Server {
	return c.server
}

func (c *Client) ReadMessage() {
	for {
		if c.isRunning == false {
			return
		}
		var messageContent []byte
		var err error
		if c.connType == ConnTypeTcp {

			reader := bufio.NewReader(c.tcpConn)
			header, err0 := reader.Peek(HeaderFlagLen + HeaderBodyLen)
			if err0 != nil {
				if err0 == io.EOF {
					continue
				} else {
					c.Close(err0)
					return
				}
			}
			if !bytes.HasPrefix(header, HeaderFlag) {
				c.Close(errors.New("header flag error"))
				return
			}
			headerBody := bytes.TrimPrefix(header, HeaderFlag)
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
			if int32(reader.Buffered()) < HeaderFlagLen+HeaderBodyLen+bodyLen {
				continue
			}
			data := make([]byte, HeaderFlagLen+HeaderBodyLen+bodyLen)
			_, err3 := reader.Read(data)
			if err3 != nil {
				if err3 == io.EOF {
					continue
				} else {
					c.Close(err3)
					return
				}
			}
			messageContent = data[HeaderFlagLen+HeaderBodyLen:]
		} else {
			_, messageContent, err = c.wsConn.ReadMessage()
			if err != nil {
				if err == io.EOF {
					continue
				} else {
					color.Red("message read error:" + err.Error())
					c.Close(err)
					return
				}
			}
		}
		color.Yellow("message received:%s", string(messageContent))
		message := NewMessage(0, c.userId, messageContent)
		c.Send(message)
	}

}
func (c *Client) WriteMessage() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		if c.isRunning == false {
			return
		}
		select {
		case <-c.closeCh:
			close(c.readCh)
			close(c.writeCh)
			ticker.Stop()
			return
		case <-ticker.C:
			message := NewMessage(0, c.userId, []byte("ping"))
			var err error
			if c.connType == ConnTypeTcp {
				content, _ := Pack("ping")
				color.Green("message send:%s", string(message.data))
				_, err = c.tcpConn.Write(content)
			} else {
				err = c.wsConn.WriteMessage(websocket.TextMessage, message.data)
			}
			if err != nil {
				c.Close(err)
				color.Red("send message error:%s", err.Error())
				return
			}
		case message, ok := <-c.writeCh:
			if ok {
				var err error
				if c.connType == ConnTypeTcp {
					content, _ := Pack(string(message.data))
					_, err = c.tcpConn.Write(content)
				} else {
					err = c.wsConn.WriteMessage(websocket.TextMessage, message.data)
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

func (c *Client) Send(message *Message) {
	if c.isRunning {
		color.Green("message send:%s", string(message.data))
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

func NewWsClient(userId uint64, clientId string, server *Server, conn *websocket.Conn) *Client {
	client := &Client{
		clientId:  clientId,
		server:    server,
		userId:    userId,
		connType:  ConnTypeWs,
		tcpConn:   nil,
		wsConn:    conn,
		readCh:    make(chan *Message, 1000),
		writeCh:   make(chan *Message, 1000),
		closeCh:   make(chan byte),
		mutex:     sync.Mutex{},
		isRunning: true,
	}
	go client.ReadMessage()
	go client.WriteMessage()
	return client
}

func NewTcpClient(userId uint64, clientId string, server *Server, conn *net.TCPConn) *Client {
	client := &Client{
		clientId:  clientId,
		server:    server,
		userId:    userId,
		connType:  ConnTypeTcp,
		tcpConn:   conn,
		wsConn:    nil,
		readCh:    make(chan *Message, 1000),
		writeCh:   make(chan *Message, 1000),
		closeCh:   make(chan byte),
		mutex:     sync.Mutex{},
		isRunning: true,
	}
	go client.ReadMessage()
	go client.WriteMessage()
	return client
}
