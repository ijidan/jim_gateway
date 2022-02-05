package manager

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

//消息类型

const (
	MessageTypeText   = iota      //text message
	MessageTypeBinary             //binary message
	MessageTypeImg                //image message
	MessageTypeVideo              //video message
	MessageTypeClose  = iota + 90 //close message
	MessageTypePing   = iota + 91 //ping message
	MessageTypePong   = iota + 92 //pong message

)

type Message struct {
	id         int64
	fromUserId uint64
	toUserId   uint64
	data       []byte
	createTime int64
}

func GenMessageId() int64 {
	node, _ := snowflake.NewNode(1)
	id := node.Generate().Int64()
	return id
}

func NewMessage(fromUserId uint64, toUserId uint64, data []byte) *Message {
	id := GenMessageId()
	message := &Message{
		id:         id,
		fromUserId: fromUserId,
		toUserId:   toUserId,
		data:       data,
		createTime: time.Now().Unix(),
	}
	return message
}
