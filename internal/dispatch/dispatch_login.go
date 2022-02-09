package dispatch

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/fatih/color"
	"jim_gateway/pkg"
)

func PublishCmdLogin(gatewayId uint64,userId uint64) error {
	kafka := pkg.NewKafkaS(pkg.Conf.PubSub.Brokers)
	cmdLogin:=CmdLogin{
		Cmd:       "cmd.login",
		GatewayId: gatewayId,
		ClientId:  userId,
		UserId:    userId,
	}
	jsonBytes,_:=json.Marshal(cmdLogin)
	return kafka.PublishMessage(TopicCmdLogin, jsonBytes)
}

func SubscribeCmdLogin() error {
	f := func(msg *sarama.ConsumerMessage) error {
		color.Yellow("message received:%s", string(msg.Value))
		return nil
	}
	kafka := pkg.NewKafkaS(pkg.Conf.PubSub.Brokers)
	return kafka.Subscribe(TopicCmdLogin, GroupCmdLogin, f)
}