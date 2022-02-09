package manager

import (
	"encoding/json"
	"github.com/fatih/color"
)

func ReadMessageFromGrpc() {
	sendClient := GetGatewayServiceSendMessageClient()
	go func() {
		for {
			rsp, err := sendClient.Recv()
			if err != nil {
				return
			}
			color.Yellow("receive grpc message success:%s", string(rsp.Data))
			var data json.RawMessage
			msg := ClientMessage{
				Data: &data,
			}
			if err3 := json.Unmarshal(rsp.Data, &msg); err3 != nil {
				color.Red("parse message err:%s", err3.Error())
			}
			color.Cyan("parse grpc message:%s",msg.Cmd)
			switch msg.Cmd {
			case "chat.c2c.txt":
				ParseC2CTxtMessage(data, rsp.Data)
			}
		}
	}()
}
