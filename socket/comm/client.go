package comm

import (
	"errors"
	"net"
	"sync"
)

type Client struct {
	addr    string
	conn    net.Conn
	UserId  uint32
	wait    sync.WaitGroup
	handler map[uint32]Handler
	sendbuf chan []byte
}

var id uint32

func NewClient(addr string) *Client {
	c := Client{addr: addr, UserId: id}

	id++

	c.sendbuf = make(chan []byte, 1000)
	c.handler = make(map[uint32]Handler, 100)

	return &c
}

func (s *Client) RegHandler(channel uint32, fun Handler) error {

	_, b := s.handler[channel]
	if b == true {
		return errors.New("channel has been register!")
	}

	s.handler[channel] = fun

	return nil
}

func (c *Client) Run() error {

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = conn

	c.wait.Add(2)
	go socketsend(c.UserId, c.conn, c.sendbuf, &c.wait)
	go socketrecv(c.UserId, c.conn, c.handler, &c.wait)
	c.wait.Wait()

	return nil
}

func (c *Client) Stop() error {

	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMsg(channel uint32, body []byte) error {

	var msghead msgHeader

	msghead.BodySize = uint32(len(body))
	msghead.Channel = channel
	msghead.MagicId = MAGIC_FLAG

	buftemp, err := codeMsgHeader(msghead)
	if err != nil {
		return err
	}

	buftemp = append(buftemp, body...)

	c.sendbuf <- buftemp

	return nil
}
