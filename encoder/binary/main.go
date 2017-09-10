package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// 报文序列化
func CodePacket(req interface{}) ([]byte, error) {
	iobuf := new(bytes.Buffer)
	err := binary.Write(iobuf, binary.BigEndian, req)
	if err != nil {
		return nil, err
	}
	return iobuf.Bytes(), nil
}

// 报文反序列化
func DecodePacket(buf []byte, rsp interface{}) error {
	iobuf := bytes.NewReader(buf)
	err := binary.Read(iobuf, binary.BigEndian, rsp)
	return err
}

type msgblock struct {
	A uint32
	B uint64
	C [2]int32
	D bool
}

func main() {
	var msg msgblock
	var msg2 msgblock

	msg.A = 1
	msg.B = 2
	msg.C[0] = -31
	msg.C[1] = -32
	msg.D = true

	buf, err := CodePacket(msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = DecodePacket(buf, &msg2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(msg2)

	if msg2 == msg {
		fmt.Println("binary encoder & decoder success!")
	}

}
