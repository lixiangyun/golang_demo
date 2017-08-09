package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	IP   = "192.168.0.107"
	PORT = "6060"
)

var recvmsgcnt int
var recvmsgsize int

func netstat() {

	lastrecvmsgcnt := recvmsgcnt
	lastrecvmsgsize := recvmsgsize

	for {
		time.Sleep(time.Duration(1 * time.Second))

		fmt.Printf("Speed %d cnt/s , %.3f MB/s\r\b",
			recvmsgcnt-lastrecvmsgcnt,
			float32(recvmsgsize-lastrecvmsgsize)/(1024*1024))

		lastrecvmsgcnt = recvmsgcnt
		lastrecvmsgsize = recvmsgsize
	}
}

func msgProc(conn net.Conn) {

	var buf [65535]byte

	defer conn.Close()

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				fmt.Println("close connect! ", conn.RemoteAddr())
				return
			}
		}

		recvmsgcnt++
		recvmsgsize += n
	}

}

func Server() {

	listen, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go netstat()

	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			fmt.Println(err.Error())
			continue
		}
		go msgProc(conn)
	}
}

func Client() {
	conn, err := net.Dial("tcp", IP+":"+PORT)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()

	var buf [1024 * 63]byte

	for {
		_, err2 := conn.Write(buf[0:])
		if err2 != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: <-s/-c>")
	}

	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
