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

type Handler func(uint32, uint32, []byte)

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
	UserId uint32

	sendbuf chan []byte
	handler map[uint32]Handler
	conn    net.Conn
	wait    *sync.WaitGroup
}

type Listen struct {
	id     uint32
	listen net.Listener
}

func NewListen(addr string) *Listen {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &Listen{listen: listen}
}

func (l *Listen) Accept() (*Server, error) {

	conn, err := l.listen.Accept()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	s := new(Server)

	s.sendbuf = make(chan []byte, 1000)
	s.handler = make(map[uint32]Handler, 100)
	s.conn = conn
	s.UserId = l.id
	s.wait = new(sync.WaitGroup)

	l.id++

	return s, nil
}

func (s *Server) Run() {

	s.wait.Add(2)

	go socketrecv(s.UserId, s.conn, s.handler, s.wait)
	go socketsend(s.UserId, s.conn, s.sendbuf, s.wait)

	s.Stop()
}

func (s *Server) Stop() {
	s.wait.Wait()
	s.conn.Close()
	log.Println("shutdown!")
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

func socketsend(userid uint32, conn net.Conn, sendbuf chan []byte, wait *sync.WaitGroup) {

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

func socketrecv(userid uint32, conn net.Conn, funtable map[uint32]Handler, wait *sync.WaitGroup) {

	var buf [MAX_BUF_SIZE]byte
	var totallen int

	defer wait.Done()

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

			fun, b := funtable[msghead.Channel]
			if b == true {
				fun(userid, msghead.Channel, buf[bodybegin:bodyend])
			} else {
				log.Println("this channel id not found!", msghead.Channel)
			}

			lastindex = bodyend
		}
	}
}

func (s *Server) SendMsg(channel uint32, body []byte) error {

	var msghead msgHeader

	msghead.BodySize = uint32(len(body))
	msghead.Channel = channel
	msghead.MagicId = MAGIC_FLAG

	buftemp, err := codeMsgHeader(msghead)
	if err != nil {
		return err
	}

	buftemp = append(buftemp, body...)

	s.sendbuf <- buftemp

	return nil
}

func (s *Server) RegHandler(channel uint32, fun Handler) error {

	_, b := s.handler[channel]
	if b == true {
		return errors.New("channel has been register!")
	}

	s.handler[channel] = fun

	return nil
}
