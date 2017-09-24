package comm

import (
	"net"
	"sync"
)

type Client struct {
	addr    string
	conn    net.Conn
	wait    sync.WaitGroup
	sendbuf chan byte
}

func NewClient(addr string) *Client {
	c := Client{addr: addr}
	c.sendbuf = make(chan byte, MAX_BUF_SIZE)
	return &c
}

func ClientRecv(conn net.Conn, wait *sync.WaitGroup) {

	defer wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {
		cnt, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		recvmsgcnt++
		recvmsgsize += cnt
	}
}

func ClientSend(conn net.Conn, wait *sync.WaitGroup) {

	defer wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {
		cnt, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		recvmsgcnt++
		recvmsgsize += cnt
	}
}

func (c *Client) Start() error {

	var wait sync.WaitGroup

	conn, err := net.Dial("tcp", IP+":"+PORT)
	if err != nil {
		log.Println(err.Error())
		return
	}

	wait.Add(2)

	go ClientSend(conn, &wait)
	go ClientRecv(conn, &wait)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
	}

	conn.Close()

	wait.Wait()
}

func (c *Client) Stop() error {

	err := c.conn.Close()
	if err != nil {
		return err
	}

	c.wait.Wait()
}

func (c *Client) SendMsg(channel uint32, body []byte) error {

	defer wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {

		cnt, err := conn.Write(buf[0:sendbuflen])
		if err != nil {
			log.Println(err.Error())
			return
		}

		sendmsgcnt++
		sendmsgsize += cnt
	}
}
