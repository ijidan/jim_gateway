package manager

import "github.com/fatih/color"

func ReadMessageFromGrpc() {
	sendClient := GetGatewayServiceSendMessageClient()
	clientManager := GetClientManagerInstance()
	go func() {
		for {
			rsp, err := sendClient.Recv()
			if err != nil {
				color.Red("received grpc message error:%s",err.Error())
				return
			}
			//gatewayId := rsp.GatewayId
			cmd :=rsp.Cmd
			requestId:=rsp.RequestId
			data:=rsp.GetData()
			receiverId:=rsp.ReceiverId

			color.Cyan("received grpc message:%s,%s,%s,%s",cmd,receiverId,receiverId,string(data))
			receiveClient := clientManager.GetClientByClientId(receiverId)

			if receiveClient != nil && receiveClient.isRunning {
				content, _ := BusinessPack(cmd, requestId, string(data))
				receiveClient.Send(content)
			}
		}
	}()
}
