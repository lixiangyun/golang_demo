package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	IP   = "localhost"
	PORT = "6060"
)

var recvmsgcnt int
var recvmsgsize int

var sendmsgcnt int
var sendmsgsize int

func netstat() {

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

	var buf [65535]byte

	defer conn.Close()

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		recvmsgcnt++
		recvmsgsize += n

		n, err = conn.Write(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		sendmsgcnt++
		sendmsgsize += n
	}
}

func Server() {

	listen, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Println(err.Error())
		return
	}

	go netstat()

	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			log.Println(err.Error())
			continue
		}
		go msgProc(conn)
	}
}

func ClientSend(conn net.Conn, wait *sync.WaitGroup) {

	defer wait.Done()
	var buf [4096]byte

	for {
		cnt, err := conn.Write(buf[0:])
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
	var buf [4096]byte

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

func Client() {

	var wait sync.WaitGroup

	conn, err := net.Dial("tcp", IP+":"+PORT)
	if err != nil {
		log.Println(err.Error())
		return
	}

	wait.Add(2)

	go ClientSend(conn, &wait)
	go ClientRecv(conn, &wait)
	go netstat()

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
	}

	conn.Close()

	wait.Wait()
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
