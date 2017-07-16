package ipc

import (
	"encoding/json"
	"fmt"
)

type IpcClient struct {
	conn chan string
}

func NewIpcClient(s *IpcServer) *IpcClient {
	c := s.Connect()
	return &IpcClient{c}
}

func (c *IpcClient) Call(m, p string) (resp *Response, err error) {
	req := &Request{m, p}

	var b []byte

	b, err = json.Marshal(req)

	if err != nil {
		return
	}

	fmt.Println("Requst2 : ", string(b))

	c.conn <- string(b)

	str := <-c.conn // 等待回应

	fmt.Println("Response2 : ", str)

	var resp1 Response

	err = json.Unmarshal([]byte(str), &resp1)

	resp = &resp1

	return
}
