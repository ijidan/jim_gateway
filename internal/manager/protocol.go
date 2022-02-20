package manager

import (
	"bytes"
	"encoding/binary"
)

const (
	BusinessHeaderFlagLen      = 4
	BusinessHeaderCmdLen       = 4
	BusinessHeaderRequestIdLen = 4
	BusinessHeaderContentLen   = 4
)

const (
	BusinessHeaderFlag = "jim1"
)

const (
	BusinessCmdPing        = "A001" //ping
	BusinessCmdPong        = "A002" //pong
	BusinessCmdServerClose = "A999" //close：服务端断开
	BusinessCmdClientClose = "A998" //close：客户端断开

	BusinessCmdAuthReq     = "B001" //请求认证
	BusinessCmdAuthSuccess = "B002" //认证成功
	BusinessCmdAuthFail    = "B003" //认证失败

	BusinessCmdC2C = "B101" //C2C单聊消息
	BusinessCmdC2G = "B102" //C2G群聊消息
	BusinessCmdS2C = "B103" //S2C系统推送消息
)

type BusinessHeader struct {
	headerFlag [4]byte
	cmd        [4]byte
	requestId  [4]byte
	ContentLen [4]byte
	Data       []byte
}

func BusinessPack(cmd string, requestId uint32, message string) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	//write header flag
	err = binary.Write(buf, binary.BigEndian, []byte(BusinessHeaderFlag))
	if err != nil {
		return nil, err
	}

	//write cmd
	err = binary.Write(buf, binary.BigEndian, []byte(cmd))
	if err != nil {
		return nil, err
	}

	//write request id
	_ = Uint32ToBytes(requestId)
	err = binary.Write(buf, binary.BigEndian, requestId)
	if err != nil {
		return nil, err
	}

	//write content len
	bodyLen := int32(len(message))
	err = binary.Write(buf, binary.BigEndian, bodyLen)
	if err != nil {
		return nil, err
	}
	//write data
	err = binary.Write(buf, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Uint32ToBytes(i uint32) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}
