package manager

import (
	"bytes"
	"encoding/binary"
)

const (
	HeaderFlagLen = 2
	HeaderBodyLen = 4
)

var HeaderFlag = []byte("jj")

type ProtocolHeader struct {
	headerFlag [3]byte
	ContentLen [4]byte
	Token      [32]byte
	Data       []byte
}

func Pack(message string) ([]byte, error) {
	buf := new(bytes.Buffer)
	var err error
	err = binary.Write(buf, binary.BigEndian, HeaderFlag)
	if err != nil {
		return nil, err
	}
	bodyLen := int32(len(message))
	err = binary.Write(buf, binary.BigEndian, bodyLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
