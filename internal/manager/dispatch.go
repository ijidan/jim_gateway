package manager

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	"jim_gateway/pkg"
)

const (
	TopicSendMessage = "topic.send.message"
)
const (
	GroupSendMessage="group.send.message"
)

func PublishSendMessage(message []byte) error {
	kafka := pkg.NewKafkaS(pkg.Conf.PubSub.Brokers)
	topic:=fmt.Sprintf("%s.%s",TopicSendMessage,pkg.Conf.Gateway.Id)
	return kafka.PublishMessage(topic, message)
}

func SubscribeSendMessage() error {
	f := func(message *sarama.ConsumerMessage) error {
		color.Yellow("kafka message received:%s", string(message.Value))
		var data json.RawMessage
		msg := ClientMessage{
			Data: &data,
		}
		if err3 := json.Unmarshal(message.Value, &msg); err3 != nil {
			color.Red("parse message err:%s", err3.Error())
		}
		switch msg.Cmd {
		case "chat.c2c.txt":
			ParseC2CTxtMessage(data, message.Value)
		}
		return nil
	}
	kafka := pkg.NewKafkaS(pkg.Conf.PubSub.Brokers)
	topic:=fmt.Sprintf("%s.%s",TopicSendMessage,pkg.Conf.Gateway.Id)
	return kafka.Subscribe(topic, GroupSendMessage, f)
}