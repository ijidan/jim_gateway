package manager

import (
	"context"
	"github.com/fatih/color"
	"jim_gateway/internal/jim_proto/proto_build"
	"jim_gateway/pkg"
	"sync"
)

var (
	onceGatewayServiceClient sync.Once
	gatewayServiceClient     proto_build.GatewayServiceClient
)

func GetGatewayServiceClientInstance() proto_build.GatewayServiceClient {
	onceGatewayServiceClient.Do(func() {
		commonServiceBasic := NewBasicCall(pkg.Conf.Rpc.Host, pkg.Conf.Rpc.Port)
		gatewayServiceClient = proto_build.NewGatewayServiceClient(commonServiceBasic.Conn)
	})
	return gatewayServiceClient
}

var (
	onceGatewayServiceSendMessageClient sync.Once
	gatewayServiceSendMessageClient     proto_build.GatewayService_SendMessageClient
)

func GetGatewayServiceSendMessageClient() proto_build.GatewayService_SendMessageClient{
	onceGatewayServiceSendMessageClient.Do(func() {
		gatewayServiceClient:=GetGatewayServiceClientInstance()
		gatewayServiceSendMessageClient,_=gatewayServiceClient.SendMessage(context.Background())
	})
	return gatewayServiceSendMessageClient
}

func RegisterGateway(gatewayId string)  {
	client:=GetGatewayServiceClientInstance()
	push,_:=client.SendMessage(context.TODO())
	sendMessage:=&proto_build.SendMessageRequest{
		GatewayId: gatewayId,
		Data:      []byte("ping"),
	}
	err:=push.Send(sendMessage)
	if err!=nil{
		color.Red("register gateway error:%s",err.Error())
	}
}

func RegisterClient(gatewayId string,clientId string) (*proto_build.RegisterResponse, error) {
	client:=GetGatewayServiceClientInstance()
	req:=&proto_build.RegisterRequest{
		GatewayId: gatewayId,
		ClientId:  clientId,
	}
	return client.Register(context.TODO(),req)
}

func UnRegisterClient(gatewayId string,clientId string)(*proto_build.UnRegisterResponse, error)  {
	client:=GetGatewayServiceClientInstance()
	req:=&proto_build.UnRegisterRequest{
		GatewayId: gatewayId,
		ClientId:  clientId,
	}
	return client.UnRegister(context.TODO(),req)
}
