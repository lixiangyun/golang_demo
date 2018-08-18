package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	CA_FILE     = "../crt/ca.crt"
	SERVER_CERT = "../crt/server.crt"
	SERVER_KEY  = "../crt/server.pem"
	CLIENT_CERT = "../crt/client.crt"
	CLIENT_KEY  = "../crt/client.pem"
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
)

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

func msgProc(conn net.Conn) {

	var buf [MAX_BUF_SIZE]byte

	defer conn.Close()

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		recvmsgcnt++
		recvmsgsize += n

		n, err = conn.Write(buf[0:n])
		if err != nil {
			log.Println(err.Error())
			return
		}

		sendmsgcnt++
		sendmsgsize += n
	}
}

func Server(addr string) {

	//这里读取的是根证书
	buf, err := ioutil.ReadFile(CA_FILE)
	if err != nil {
		log.Println(err.Error())
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(buf)

	//加载服务端证书
	cert, err := tls.LoadX509KeyPair(SERVER_CERT, SERVER_KEY)
	if err != nil {
		log.Println(err.Error())
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	go netstat_server()

	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			log.Println(err.Error())
			continue
		}

		//accept到conn后
		tlsConn := tls.Server(conn, tlsConfig)

		go msgProc(tlsConn)
	}
}

var sendbuflen = 128

func ClientSend(conn net.Conn, wait *sync.WaitGroup) {

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

func Client(addr string) {

	var wait sync.WaitGroup

	//这里读取的是根证书
	buf, err := ioutil.ReadFile(CA_FILE)
	if err != nil {
		log.Println(err.Error())
		return
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(buf)

	//加载客户端证书
	//这里加载的是服务端签发的
	cert, err := tls.LoadX509KeyPair(CLIENT_CERT, CLIENT_KEY)
	if err != nil {
		log.Fatalln(err)
	}

	config := &tls.Config{
		ServerName:   "127.0.0.1",
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	tlsConn := tls.Client(conn, config)

	wait.Add(2)

	go ClientSend(tlsConn, &wait)
	go ClientRecv(tlsConn, &wait)
	go netstat_client()

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
	}

	tlsConn.Close()

	wait.Wait()
}

func main() {

	args := os.Args

	if len(args) < 3 {
		log.Println("Usage: <-s/-c> <ip:port>")
		return
	}

	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2])
	}
}
