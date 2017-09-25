package comm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"os"
	"sync"
)

type Handler func(uint32, []byte)

const (
	MAX_BUF_SIZE = 128 * 1024
	MAGIC_FLAG   = 0x98b7f30a
	MSG_HEAD_LEN = 4 * 4
)

type msgHeader struct {
	MagicId   uint32
	Channel   uint32
	RequestId uint32
	BodySize  uint32
}

func FullyWrite(conn net.Conn, buf []byte) error {

	totallen := len(buf)
	sendcnt := 0

	for {

		cnt, err := conn.Write(buf[sendcnt:])
		if err != nil {
			return err
		}

		if cnt <= 0 {
			return errors.New("conn write error!")
		}

		if cnt+sendcnt >= totallen {
			return nil
		}

		sendcnt += cnt
	}
}

type Server struct {
	addr      string
	requestid uint32
	sendbuf   chan []byte
	handler   map[uint32]Handler
	listen    net.Listener
	lock      sync.RWMutex
	wait      sync.WaitGroup
}

func NewServer(addr string) *Server {

	s := Server{addr: addr}

	s.sendbuf = make(chan []byte, 1000)
	s.handler = make(map[uint32]Handler, 100)

	return &s
}

func (s *Server) Start() error {

	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.wait.Add(1)

	go func() {

		defer s.wait.Done()

		for {

			conn, err := listen.Accept()
			if err != nil {
				log.Println(err.Error())
				return
			}

			s.wait.Add(2)
			go socketrecv(conn, s.handler[0], &s.wait)
			go socketsend(conn, s.sendbuf, &s.wait)
		}
	}()

	return nil
}

func (s *Server) Stop() error {

	err := s.listen.Close()
	if err != nil {
		return err
	}

	s.wait.Wait()

	return nil
}

func codeMsgHeader(req msgHeader) ([]byte, error) {
	iobuf := new(bytes.Buffer)

	err := binary.Write(iobuf, binary.BigEndian, req)
	if err != nil {
		return nil, err
	}

	return iobuf.Bytes(), nil
}

func decodeMsgHeader(buf []byte) (rsp msgHeader, err error) {

	iobuf := bytes.NewReader(buf)
	err = binary.Read(iobuf, binary.BigEndian, &rsp)

	return
}

func socketsend(conn net.Conn, sendbuf chan []byte, wait *sync.WaitGroup) {

	defer wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {

		var buflen int

		tmpbuf := <-sendbuf
		tmpbuflen := len(tmpbuf)

		if tmpbuflen >= MAX_BUF_SIZE/2 {
			err := FullyWrite(conn, tmpbuf[0:])
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			copy(buf[0:tmpbuflen], tmpbuf[0:])
			buflen = tmpbuflen
		}

		chanlen := len(sendbuf)

		for i := 0; i < chanlen; i++ {

			tmpbuf = <-sendbuf
			tmpbuflen = len(tmpbuf)

			copy(buf[buflen:buflen+tmpbuflen], tmpbuf[0:])
			buflen += tmpbuflen

			if buflen >= MAX_BUF_SIZE/2 {
				err := FullyWrite(conn, buf[0:buflen])
				if err != nil {
					log.Println(err.Error())
					return
				}
				buflen = 0
			}
		}

		if buflen > 0 {
			err := FullyWrite(conn, buf[0:buflen])
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func socketrecv(conn net.Conn, fun Handler, wait *sync.WaitGroup) {

	var buf [MAX_BUF_SIZE]byte
	var totallen int

	defer wait.Done()
	defer conn.Close()

	for {

		var lastindex int

		recvnum, err := conn.Read(buf[totallen:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		totallen += recvnum

		for {

			if lastindex+MSG_HEAD_LEN > totallen {
				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				totallen = totallen - lastindex
				break
			}

			msghead, err2 := decodeMsgHeader(buf[lastindex : lastindex+MSG_HEAD_LEN])
			if err2 != nil {
				log.Println(err2.Error())
				break
			}

			bodybegin := lastindex + MSG_HEAD_LEN
			bodyend := bodybegin + int(msghead.BodySize)

			if msghead.MagicId != MAGIC_FLAG {

				log.Println("msghead_0:", msghead)
				log.Println("totallen:", totallen)
				log.Println("bodybegin:", bodybegin, " bodyend:", bodyend)
				log.Println("body:", buf[lastindex:bodyend])
				log.Println("bodyFull:", buf[0:totallen])
				os.Exit(11)
			}

			if bodyend > totallen {

				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				totallen = totallen - lastindex
				break
			}

			fun(msghead.Channel, buf[bodybegin:bodyend])

			lastindex = bodyend
		}
	}
}

func (s *Server) SendMsg(channel uint32, body []byte) error {

	var msghead msgHeader

	msghead.BodySize = uint32(len(body))
	msghead.Channel = channel
	msghead.MagicId = MAGIC_FLAG
	msghead.RequestId = s.requestid

	s.requestid++

	buftemp, err := codeMsgHeader(msghead)
	if err != nil {
		return err
	}

	buftemp = append(buftemp, body...)

	s.sendbuf <- buftemp

	return nil
}

func (s *Server) RegHandler(channel uint32, fun Handler) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	_, b := s.handler[channel]
	if b == true {
		return errors.New("channel has been register!")
	}

	s.handler[channel] = fun

	return nil
}
