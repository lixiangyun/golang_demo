package comm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var recvmsgcnt int
var recvmsgsize int

var sendmsgcnt int
var sendmsgsize int

type banchmark struct {
	sendbuflen  int
	recvmsgsize int
}

const (
	MAX_BUF_SIZE = 128 * 1024
	MAGIC_FLAG   = 0x98b7f30a
	MSG_HEAD_LEN = 4 * 4
)

type msgHeader struct {
	magicid   uint32
	channel   uint32
	requestid uint32
	bodysize  uint32
}

var banchmarktest [20]banchmark

func netstat_client() {

	num := 0

	time.Sleep(time.Second)

	log.Println("banch mark test start...")

	lastrecvmsgcnt := recvmsgcnt
	lastrecvmsgsize := recvmsgsize

	lastsendmsgcnt := sendmsgcnt
	lastsendmsgsize := sendmsgsize

	for {

		time.Sleep(time.Second)

		log.Printf("Recv %d cnt/s , %.3f MB/s \r\n",
			recvmsgcnt-lastrecvmsgcnt,
			float32(recvmsgsize-lastrecvmsgsize)/(1024*1024))

		log.Printf("Send %d cnt/s , %.3f MB/s \r\n",
			sendmsgcnt-lastsendmsgcnt,
			float32(sendmsgsize-lastsendmsgsize)/(1024*1024))

		banchmarktest[num].sendbuflen = sendbuflen
		banchmarktest[num].recvmsgsize = recvmsgsize - lastrecvmsgsize

		num++

		lastrecvmsgcnt = recvmsgcnt
		lastrecvmsgsize = recvmsgsize

		lastsendmsgcnt = sendmsgcnt
		lastsendmsgsize = sendmsgsize

		if sendbuflen*2 <= MAX_BUF_SIZE {
			sendbuflen = sendbuflen * 2
		}

		if num >= len(banchmarktest) {
			log.Println("banch mark test end.")
			break
		}
	}

	for _, v := range banchmarktest {

		log.Printf("SendBufLen %d , %.3f MB/s \r\n",
			v.sendbuflen, float32(v.recvmsgsize)/(1024*1024))
	}
}

func netstat_server() {

	lastrecvmsgcnt := recvmsgcnt
	lastrecvmsgsize := recvmsgsize

	lastsendmsgcnt := sendmsgcnt
	lastsendmsgsize := sendmsgsize

	for {

		time.Sleep(time.Second)

		log.Printf("Recv %d cnt/s , %.3f MB/s \r\n",
			recvmsgcnt-lastrecvmsgcnt,
			float32(recvmsgsize-lastrecvmsgsize)/(1024*1024))

		log.Printf("Send %d cnt/s , %.3f MB/s \r\n",
			sendmsgcnt-lastsendmsgcnt,
			float32(sendmsgsize-lastsendmsgsize)/(1024*1024))

		lastrecvmsgcnt = recvmsgcnt
		lastrecvmsgsize = recvmsgsize

		lastsendmsgcnt = sendmsgcnt
		lastsendmsgsize = sendmsgsize
	}
}

type Handler interface {
	ChanHandler(uint32, []byte)
}

type Server struct {
	addr    string
	handler map[uint32]Handler
	listen  net.Listener
	lock    sync.RWMutex
	wait    sync.WaitGroup
}

func NewServer(addr string) *Server {

	s := Server{addr: addr}

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

			s.wait.Add(1)
			go msgrecvtask(conn, s)
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

func gethandler(s *Server, channel uint32) Handler {
	s.lock.RLock()
	fun, b := s.handler[channel]
	s.lock.RUnlock()

	if b == true {
		return fun
	} else {
		return nil
	}
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

func msgrecvtask(conn net.Conn, s *Server) {

	var buf [MAX_BUF_SIZE]byte
	var lastindex uint32
	var totallen uint32
	var nextsize uint32
	var msghead msgHeader

	defer conn.Close()

	for {

		recvnum, err := conn.Read(buf[lastindex:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		totallen = lastindex + uint32(recvnum)

		for {

			msghead, err2 := decodeMsgHeader(buf[lastindex : lastindex+MSG_HEAD_LEN])
			if err2 != nil {
				log.Println(err2.Error())
				break
			}

			bodybegin := lastindex + MSG_HEAD_LEN
			bodyend := bodybegin + msghead.bodysize

			if bodyend > totallen {
				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				lastindex = 0
				break
			}

			fun := gethandler(s, msghead.channel)
			if fun == nil {
				log.Println("can not found channel handler!", msghead)
				break
			}

			fun.ChanHandler(msghead.channel, buf[bodybegin:bodyend])

			lastindex = bodyend
		}
	}
}

func (s *Server) SendMsg(channel uint32, body []byte) error {

	return
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
