package main

import (
	"flag"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ADDRESS      string
	PARALLEL_NUM int
	RUNTIME      int
	BODY_LENGTH  int
	ROLE         string
	h            bool
)

var gStat *Stat

func init() {
	flag.StringVar(&ROLE, "r", "s", "the tools role (s/c).")
	flag.IntVar(&PARALLEL_NUM, "p", 1, "parallel tcp connect.")
	flag.IntVar(&RUNTIME, "t", 30, "total run time (second).")
	flag.IntVar(&BODY_LENGTH, "l", 64, "transport body length (KB).")
	flag.StringVar(&ADDRESS, "b", "127.0.0.1:8010", "set the service address.")
	flag.BoolVar(&h, "h", false, "this help.")
}

func ServerProc(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, BODY_LENGTH)

	for {
		cnt, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		gStat.Add(cnt, 0)

		cnt, err = conn.Write(buf[0:cnt])
		if err != nil {
			log.Println(err.Error())
			break
		}
	}
}

func Server(addr string) {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		conn, err2 := listen.Accept()
		if err2 != nil {
			log.Println(err.Error())
			continue
		}
		go ServerProc(conn)
	}
}

func ClientSend(conn net.Conn, wait *sync.WaitGroup) {
	defer wait.Done()
	buf := make([]byte, BODY_LENGTH)

	for {
		_, err := conn.Write(buf[:])
		if err != nil {
			return
		}
	}
}

func ClientRecv(conn net.Conn, wait *sync.WaitGroup) {
	defer wait.Done()
	buf := make([]byte, BODY_LENGTH)

	for {
		cnt, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		gStat.Add(cnt, 0)
	}
}

func ClientConn(addr string, client *sync.WaitGroup) {
	defer client.Done()
	var wait sync.WaitGroup

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	wait.Add(2)

	go ClientSend(conn, &wait)
	go ClientRecv(conn, &wait)

	time.Sleep(time.Duration(RUNTIME) * time.Second)
	conn.Close()

	wait.Wait()
}

func Client(addr string) {
	var wait sync.WaitGroup
	wait.Add(PARALLEL_NUM)
	for i := 0; i < PARALLEL_NUM; i++ {
		go ClientConn(addr, &wait)
	}
	wait.Wait()
}

func main() {

	flag.Parse()
	if h || (ROLE != "s" && ROLE != "c") {
		flag.Usage()
		return
	}
	BODY_LENGTH = BODY_LENGTH * 1024

	gStat = NewStat(5)

	switch ROLE {
	case "s":
		gStat.Prefix("tcp server")
		Server(ADDRESS)
	case "c":
		gStat.Prefix("tcp client")
		Client(ADDRESS)
	}

	gStat.Delete()
}
