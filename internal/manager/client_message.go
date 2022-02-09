package manager

import (
	"encoding/json"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"regexp"
)

type ClientMessage struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

type TextMessage struct {
	Id           uint64 `json:"id"`
	Content      string `json:"content"`
	ToReceiverId string `json:"to_receiver_id"`
	AtUserId     string `json:"at_user_id"`
}

type AuthMessage struct {
	Token string `json:"token"`
}

//type LocationMessage struct {
//	Id uint64 `json:"id"`
//	cover_image string `json:"cover_image"`
//	Float lat
//	double lng = 4;
//	string map_link = 5;
//	string desc = 6;
//}
//
//
////表情消息
//message FaceMessage{
//uint64 id = 1;
//string symbol = 2;
//}
//*/

//
//type C2CMessage struct {
//	Type    string `yaml:"type"`
//	Content string `yaml:"content"`
//}
//
//type C2GMessage struct {
//	Type    string `yaml:"type"`
//	Content string `yaml:"content"`
//}

//func ParseMessage(mode string, messageContent []byte) (string, uint64, error) {
//	var clientId string
//	var userId uint64
//	var err error
//	if mode == ModeLocal.String() {
//		clientId, userId, err = ParseTcpMessage(messageContent)
//	}
//	if mode == ModelGrpc.String() {
//		clientId, userId, err = ParseGrpcMessage(messageContent)
//	}
//	if mode == ModelKafka.String() {
//	}
//	return clientId, userId, err
//}

func ParseAuthReqMessage(data json.RawMessage)(string,uint64)  {
	var content AuthMessage
	if err4 := json.Unmarshal(data, &content); err4 != nil {
		color.Red("parse message err:%s", err4.Error())
	}
	token:=content.Token
	clientId:=token

	re := regexp.MustCompile("[0-9]+")
	all:=re.FindAllString(clientId, -1)
	userId:=cast.ToUint64(all[0])

	return clientId,userId
}

func ParseC2CTxtMessage(data json.RawMessage,messageContent []byte)  {
	clientManager := GetClientManagerInstance()
	var content TextMessage
	if err4 := json.Unmarshal(data, &content); err4 != nil {
		color.Red("parse message err:%s", err4.Error())
	}

	toReceiverId := content.ToReceiverId
	color.Cyan("parse txt message:%s",toReceiverId)
	toClient := clientManager.GetClientByClientId(toReceiverId)

	if toClient != nil && toClient.isRunning {
		toClient.Send(messageContent)
	}
}

//func ParseGrpcMessage(messageContent []byte) (string, uint64, error) {
//	req := &proto_build.SendMessageRequest{
//		GatewayId: 1,
//		Data:      messageContent,
//	}
//	sendClient := call.GetGatewayServiceSendMessageClient()
//	errSend1 := sendClient.Send(req)
//	if errSend1 != nil {
//		color.Red("send client send message error:%s", errSend1.Error())
//	}
//	return "", 0, nil
//}
//
//func ParseTcpMessage(messageContent []byte) (string, uint64, error) {
//	clientManager := GetClientManagerInstance()
//	var data json.RawMessage
//	msg := ClientMessage{
//		Data: &data,
//	}
//	if err3 := json.Unmarshal(messageContent, &msg); err3 != nil {
//		color.Red("parse message err:%s", err3.Error())
//	}
//	switch msg.Cmd {
//	case "auth.req":
//		var content AuthMessage
//		if err4 := json.Unmarshal(data, &content); err4 != nil {
//			color.Red("parse message err:%s", err4.Error())
//		}
//		token := content.Token
//		clientId := token
//
//		re := regexp.MustCompile("[0-9]+")
//		all := re.FindAllString(clientId, -1)
//		userId := all[0]
//
//		return clientId, cast.ToUint64(userId), nil
//	case "chat.c2c.txt":
//		var content TextMessage
//		if err4 := json.Unmarshal(data, &content); err4 != nil {
//			color.Red("parse message err:%s", err4.Error())
//		}
//		toReceiverId := content.ToReceiverId
//		toClient := clientManager.GetClientByClientId(toReceiverId)
//		if toClient != nil && toClient.isRunning {
//			toClient.Send(messageContent)
//		}
//	}
//	return "", 0, nil
//}
