package main

import (
	"golang_demo/socket/comm"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	IP   = "localhost"
	PORT = "1234"
)

var flag chan int

var servertable map[uint32]*comm.Server
var clienttable map[uint32]*comm.Client

func serverhandler(userid uint32, channel uint32, body []byte) {

	err := servertable[userid].SendMsg(channel, body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	recvmsgcnt++
	recvmsgsize += len(body)

	sendmsgcnt++
	sendmsgsize += len(body)
}

var recvmsgcnt int
var recvmsgsize int

var sendmsgcnt int
var sendmsgsize int

var sendbuflen = 128

type banchmark struct {
	sendbuflen  int
	recvmsgsize int
}

var banchmarktest [100]banchmark

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

		if num >= len(banchmarktest) {
			log.Println("banch mark test end.")
			break
		}
	}

	for _, v := range banchmarktest {

		log.Printf("SendBufLen %d , %.3f MB/s \r\n",
			v.sendbuflen, float32(v.recvmsgsize)/(1024*1024))
	}

	flag <- 0
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

func Server() {

	servertable = make(map[uint32]*comm.Server, 1000)

	list := comm.NewListen(":" + PORT)

	go netstat_server()

	var server *comm.Server
	var err error

	for {
		server, err = list.Accept()
		if err != nil {
			log.Println(err.Error())
			return
		}
		server.RegHandler(0, serverhandler)

		servertable[server.UserId] = server

		go server.Run()
	}
}

func clienthandler(userid uint32, channel uint32, body []byte) {
	recvmsgcnt++
	recvmsgsize += len(body)
}

func Client() {

	clienttable = make(map[uint32]*comm.Client, 1000)

	flag = make(chan int)

	c := comm.NewClient(IP + ":" + PORT)

	c.RegHandler(0, clienthandler)

	clienttable[c.UserId] = c

	go c.Run()
	go netstat_client()

	var sendbuf [comm.MAX_BUF_SIZE]byte

	for {
		err := c.SendMsg(0, sendbuf[0:sendbuflen])
		if err != nil {
			log.Println(err.Error())
			return
		}
		sendmsgcnt++
		sendmsgsize += sendbuflen
	}

	<-flag

	c.Stop()
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	args := os.Args

	if len(args) < 2 {
		log.Println("Usage: <-s/-c>")
	}

	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
