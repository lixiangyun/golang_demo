package ipc

import (
	"encoding/json"
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

	c.conn <- string(b)

	str := <-c.conn // 等待回应

	var resp1 Response

	err = json.Unmarshal([]byte(str), &resp1)

	resp = &resp1

	return
}
