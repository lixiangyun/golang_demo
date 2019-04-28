package main

import (
	"encoding/binary"
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
		readcnt, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		gStat.Add(readcnt, 0)

		var sendcnt int
		for {
			cnt, err := conn.Write(buf[sendcnt:readcnt])
			if err != nil {
				log.Println(err.Error())
				return
			}
			sendcnt += cnt
			if sendcnt >= readcnt {
				break
			}
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
	var sendidx uint64

	buf := make([]byte, BODY_LENGTH)
	num := BODY_LENGTH / 8

	for {

		for i := 0; i < num; i++ {
			sendidx++
			binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], sendidx)
		}

		var sendcnt int
		for {
			cnt, err := conn.Write(buf[:])
			if err != nil {
				return
			}
			sendcnt += cnt
			if sendcnt >= len(buf) {
				break
			}
		}
	}
}

func ClientRecv(conn net.Conn, wait *sync.WaitGroup) {
	defer wait.Done()

	var sendidx uint64
	buf := make([]byte, BODY_LENGTH)

	var remain int

	for {
		cnt, err := conn.Read(buf[remain:])
		if err != nil {
			return
		}
		remain += cnt

		if remain%8 == 0 {
			num := remain / 8
			for i := 0; i < num; i++ {
				idx := binary.BigEndian.Uint64(buf[i*8 : (i+1)*8])
				if idx != sendidx+1 {
					log.Fatalln("recv err body data!", idx, sendidx)
				}
				sendidx = idx
			}

			gStat.Add(remain, 0)
			remain = 0
		}
	}
}

func ClientTimer(conn net.Conn) {
	time.Sleep(time.Duration(RUNTIME) * time.Second)
	conn.Close()
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
	go ClientTimer(conn)

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
