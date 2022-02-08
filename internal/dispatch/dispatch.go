package dispatch

const (
	TopicCmdLogin = "topic.cmd.login"
)
const (
	GroupCmdLogin="group.cmd.login"
)

type CmdLogin struct {
	Cmd       string `json:"cmd"`
	GatewayId uint64 `json:"gateway_id"`
	ClientId  uint64 `yaml:"clientId"`
	UserId    uint64 `yaml:"user_id"`
}

type CmdC2C struct {
	Cmd string `json:"cmd"`
}
