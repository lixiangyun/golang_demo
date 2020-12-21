package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func Server(port string) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err2 := net.ListenUDP("udp", addr)
	if err2 != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	var buf [65535]byte
	for {
		cnt, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			if err == io.EOF {
				fmt.Println("close connect! ", conn.RemoteAddr())
				return
			}
		}
		conn.WriteToUDP(buf[:cnt], addr)
	}
}

func Client(addr string, local string) {
	localaddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	remoteaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn, err := net.ListenUDP("udp", localaddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	fmt.Println( "server : " + remoteaddr.String())

	go func() {
		for {
			writebuff := []byte("hello world!")
			cnt, err := conn.WriteToUDP(writebuff[:], remoteaddr)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if cnt > 0 {
				fmt.Printf("send to server %s success\n", remoteaddr.String())
			}
			time.Sleep(time.Millisecond)
		}
	}()

	var readbuff [1024]byte

	for  {
		cnt, addr, err := conn.ReadFromUDP(readbuff[:])
		fmt.Printf("client from: %s %s\n", addr.String(), string(readbuff[:cnt]))

		if err != nil {
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
		Client(args[2],args[3])
	}
}
