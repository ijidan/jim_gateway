package manager

type GatewayRuntimeMode string

func (r GatewayRuntimeMode) String() string {
	return string(r)
}

const (
	ModeLocal  GatewayRuntimeMode = "local"
	ModelGrpc  GatewayRuntimeMode = "grpc"
	ModelKafka GatewayRuntimeMode = "kafka"
)
