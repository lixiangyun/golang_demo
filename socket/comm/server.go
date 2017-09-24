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
	MagicId   uint32
	Channel   uint32
	RequestId uint32
	BodySize  uint32
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

		num++

		lastrecvmsgcnt = recvmsgcnt
		lastrecvmsgsize = recvmsgsize

		lastsendmsgcnt = sendmsgcnt
		lastsendmsgsize = sendmsgsize

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
	handler   map[uint32]func(uint32, []byte)
	listen    net.Listener
	lock      sync.RWMutex
	wait      sync.WaitGroup
}

func NewServer(addr string) *Server {

	s := Server{addr: addr}

	s.sendbuf = make(chan []byte, 1000)

	s.handler = make(map[uint32]func(uint32, []byte), 100)

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
			go msgsendtask(conn, s.sendbuf, s)
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

func gethandler(s *Server, channel uint32) func(uint32, []byte) {
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

	//log.Println("buf : ", buf)

	iobuf := bytes.NewReader(buf)
	err = binary.Read(iobuf, binary.BigEndian, &rsp)

	return
}

func msgsendtask(conn net.Conn, sendbuf chan []byte, s *Server) {

	defer s.wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {

		buflen := 0

		tmpbuf := <-sendbuf
		tmpbuflen := len(tmpbuf)

		copy(buf[buflen:buflen+tmpbuflen], tmpbuf[0:])
		buflen += tmpbuflen

		chanlen := len(sendbuf)

		for i := 0; i < chanlen; i++ {

			tmpbuf = <-sendbuf
			tmpbuflen = len(tmpbuf)

			copy(buf[buflen:buflen+tmpbuflen], tmpbuf[0:])
			buflen += tmpbuflen

			if buflen > MAX_BUF_SIZE/2 {
				break
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

func msgrecvtask(conn net.Conn, s *Server) {

	var buf [MAX_BUF_SIZE]byte
	var lastindex int
	var totallen int

	defer conn.Close()

	for {

		recvnum, err := conn.Read(buf[lastindex:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		//log.Println("lastindex:", lastindex)
		//log.Println("recvnum:", recvnum)
		//log.Println("body:", buf[lastindex:lastindex+recvnum])

		totallen = lastindex + recvnum

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
			bodyend := bodybegin + int(msghead.BodySize)

			//log.Println("msghead:", msghead)

			if bodyend > totallen {
				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				lastindex = 0
				break
			}

			fun := gethandler(s, msghead.Channel)
			if fun == nil {
				log.Println("can not found channel handler!", msghead)
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

func (s *Server) RegHandler(channel uint32, fun func(uint32, []byte)) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	_, b := s.handler[channel]
	if b == true {
		return errors.New("channel has been register!")
	}

	s.handler[channel] = fun

	return nil
}
