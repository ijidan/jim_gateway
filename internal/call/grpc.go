package call

import (
	"github.com/fatih/color"
)

func ReadMessageFromGrpc(){
	sendClient := GetGatewayServiceSendMessageClient()
	go func() {
		for {
			rsp, err := sendClient.Recv()
			if err != nil {
				return 
			}
			color.Yellow("receive grpc message success:%s", string(rsp.Data))
		}
	}()
}
