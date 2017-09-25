package comm

import (
	"net"
	"sync"
)

type Client struct {
	addr      string
	conn      net.Conn
	requestid uint32
	wait      sync.WaitGroup
	handler   Handler
	sendbuf   chan []byte
}

func NewClient(addr string, handler Handler) *Client {
	c := Client{addr: addr}
	c.sendbuf = make(chan []byte, 1000)
	c.handler = handler
	return &c
}

func (c *Client) Start() error {

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = conn

	c.wait.Add(2)
	go socketsend(c.conn, c.sendbuf, &c.wait)
	go socketrecv(c.conn, c.handler, &c.wait)

	return nil
}

func (c *Client) Stop() error {

	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.wait.Wait()

	return nil
}

func (c *Client) SendMsg(channel uint32, body []byte) error {

	var msghead msgHeader

	msghead.BodySize = uint32(len(body))
	msghead.Channel = channel
	msghead.MagicId = MAGIC_FLAG
	msghead.RequestId = c.requestid

	c.requestid++

	buftemp, err := codeMsgHeader(msghead)
	if err != nil {
		return err
	}

	buftemp = append(buftemp, body...)

	c.sendbuf <- buftemp

	return nil
}
