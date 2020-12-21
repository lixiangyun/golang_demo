package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	LOCAL_ADDR   string
	REMOTE_ADDR  string

	PARALLEL_NUM int
	RUNTIME      int
	BODY_LENGTH  int
	ROLE         string
	h            bool
)

var gStat *Stat

func init() {
	flag.StringVar(&ROLE, "role", "s", "the tools role (s/c).")
	flag.IntVar(&PARALLEL_NUM, "par", 1, "parallel tcp connect.")
	flag.IntVar(&RUNTIME, "runtime", 30, "total run time (second).")
	flag.IntVar(&BODY_LENGTH, "body", 64, "transport body length (KB).")
	flag.StringVar(&LOCAL_ADDR, "local", "0.0.0.0:8001", "local listem address.")
	flag.StringVar(&REMOTE_ADDR, "remote", "0.0.0.0:8001", "remote service address.")
	flag.BoolVar(&h, "help", false, "this help.")
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

func ClientConn(addr string, remote string, client *sync.WaitGroup) {
	defer client.Done()
	var wait sync.WaitGroup

	list := strings.Split(addr,":")
	addr = fmt.Sprintf("%s:%d", list[0], rand.Int()%50000 )

	localAdd, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	remoteAdd, err := net.ResolveTCPAddr("tcp", remote)
	if err != nil {
		log.Println(err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", localAdd, remoteAdd)
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

func Client(addr, remote string) {
	var wait sync.WaitGroup
	wait.Add(PARALLEL_NUM)
	for i := 0; i < PARALLEL_NUM; i++ {
		go ClientConn(addr, remote, &wait)
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
		Server(LOCAL_ADDR)
	case "c":
		gStat.Prefix("tcp client")
		Client(LOCAL_ADDR, REMOTE_ADDR)
	}

	gStat.Delete()
}
