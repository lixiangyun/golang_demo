package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
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

func msgProc(conn *net.UDPConn) {

	var buf [65535]byte

	for {
		n, _, err := conn.ReadFromUDP(buf[0:])
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

func Server(port string) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(addr)

	go netstat()

	conn, err2 := net.ListenUDP("udp", addr)
	if err2 != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()

	for {
		msgProc(conn)
	}
}

func Client(addr string) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close()

	var buf [1024]byte

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

	if len(args) < 3 {
		fmt.Println("Usage: <-s/-c> <ip:port>")
		return
	}

	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2])
	}
}
