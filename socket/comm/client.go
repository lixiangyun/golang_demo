package comm

import (
	"log"
	"net"
	"sync"
)

type Client struct {
	addr      string
	conn      net.Conn
	requestid uint32
	wait      sync.WaitGroup
	handler   func(uint32, []byte)
	sendbuf   chan []byte
}

func NewClient(addr string, handler func(uint32, []byte)) *Client {
	c := Client{addr: addr}
	c.sendbuf = make(chan []byte, 1000)
	c.handler = handler
	return &c
}

func (c *Client) ClientRecv() {

	defer c.wait.Done()
	var buf [MAX_BUF_SIZE]byte
	var lastindex uint32
	var totallen uint32

	for {

		recvnum, err := c.conn.Read(buf[lastindex:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		totallen = lastindex + uint32(recvnum)

		for {

			if lastindex+MSG_HEAD_LEN > totallen {
				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				lastindex = 0
				break
			}

			msghead, err2 := decodeMsgHeader(buf[lastindex : lastindex+MSG_HEAD_LEN])
			if err2 != nil {
				log.Println(err2.Error())
				break
			}

			bodybegin := lastindex + MSG_HEAD_LEN
			bodyend := bodybegin + msghead.BodySize

			if bodyend > totallen {
				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				lastindex = 0
				break
			}

			c.handler(msghead.Channel, buf[bodybegin:bodyend])

			lastindex = bodyend
		}
	}
}

func (c *Client) ClientSend() {

	defer c.wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {

		buflen := 0

		tmpbuf := <-c.sendbuf
		tmpbuflen := len(tmpbuf)

		copy(buf[buflen:buflen+tmpbuflen], tmpbuf[0:])
		buflen += tmpbuflen

		chanlen := len(c.sendbuf)

		for i := 0; i < chanlen; i++ {

			tmpbuf = <-c.sendbuf
			tmpbuflen = len(tmpbuf)

			copy(buf[buflen:buflen+tmpbuflen], tmpbuf[0:])
			buflen += tmpbuflen

			if buflen > MAX_BUF_SIZE/2 {
				break
			}
		}

		if buflen > 0 {
			err := FullyWrite(c.conn, buf[0:buflen])
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func (c *Client) Start() error {

	var wait sync.WaitGroup

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = conn

	wait.Add(2)

	go c.ClientSend()
	go c.ClientRecv()

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
