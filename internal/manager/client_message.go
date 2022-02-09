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

type AuthMessage struct {
	Token string `json:"token"`
}


type TextMessage struct {
	Id           uint64 `json:"id"`
	Content      string `json:"content"`
	ToReceiverId string `json:"to_receiver_id"`
	AtUserId     string `json:"at_user_id"`
}

type LocationMessage struct {
	Id         uint64  `json:"id"`
	CoverImage string  `json:"cover_image"`
	Lat        float32 `json:"lat"`
	Lng        float32 `json:"lng"`
	MapLink    string  `json:"map_link"`
	Desc       string  `json:"desc"`
}

type FaceMessage struct {
	Id     uint64 `json:"id"`
	Symbol string `json:"symbol"`
}

type SoundMessage struct {
	Id      uint64 `json:"id"`
	Url     string `json:"url"`
	Size    uint64 `json:"size"`
	Seconds uint64 `json:"seconds"`
}
type ImageMessageItem struct {
	Type   uint64 `json:"type"`
	Format uint64 `json:"format"`
	Size   uint64 `json:"size"`
	Width  uint64 `json:"width"`
	Height uint64 `json:"height"`
	Url    string `json:"url"`
}

type ImageMessage struct {
	Id               uint64             `json:"id"`
	ImageMessageItem []ImageMessageItem `json:"image_message_item"`
}

type FileMessage struct {
	Id   uint64 `json:"id"`
	Size uint64 `json:"size"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type VideoMessage struct {
	Id          uint64 `json:"id"`
	Size        uint64 `json:"size"`
	Seconds     uint64 `json:"seconds"`
	Url         string `json:"url"`
	Format      string `json:"format"`
	ThumbUrl    string `json:"thumb_url"`
	ThumbSize   uint64 `json:"thumb_size"`
	ThumbWidth  uint64 `json:"thumb_width"`
	ThumbHeight uint64 `json:"thumb_height"`
	ThumbFormat uint64 `json:"thumb_format"`
}

func ParseAuthReqMessage(data json.RawMessage) (string, uint64) {
	var content AuthMessage
	if err4 := json.Unmarshal(data, &content); err4 != nil {
		color.Red("parse message err:%s", err4.Error())
	}
	token := content.Token
	clientId := token

	re := regexp.MustCompile("[0-9]+")
	all := re.FindAllString(clientId, -1)
	userId := cast.ToUint64(all[0])

	return clientId, userId
}

func ParseC2CTxtMessage(data json.RawMessage, messageContent []byte) {
	clientManager := GetClientManagerInstance()
	var content TextMessage
	if err4 := json.Unmarshal(data, &content); err4 != nil {
		color.Red("parse message err:%s", err4.Error())
	}

	toReceiverId := content.ToReceiverId
	color.Cyan("parse txt message:%s", toReceiverId)
	toClient := clientManager.GetClientByClientId(toReceiverId)

	if toClient != nil && toClient.isRunning {
		toClient.Send(messageContent)
	}
}
