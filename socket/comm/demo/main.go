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
	PORT = "6060"
)

var flag chan int

var s *comm.Server
var c *comm.Client

func serverhandler(channel uint32, body []byte) {
	err := s.SendMsg(channel, body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//log.Println("send buf to client!", channel, body)
}

func Server() {

	s = comm.NewServer(":" + PORT)

	err := s.RegHandler(0, serverhandler)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = s.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {
		time.Sleep(time.Second)
	}

	s.Stop()
}

var clientsendcnt int
var clientrecvcnt int

func clienthandler(channel uint32, body []byte) {
	if len(body) == 16 {
		clientrecvcnt++
	}

	if clientrecvcnt >= 100000 {
		flag <- 0

		log.Println("recv buf from server!", clientrecvcnt, body)
	}
}

func Client() {

	flag = make(chan int)

	c = comm.NewClient(IP+":"+PORT, clienthandler)

	err := c.Start()
	if err != nil {
		log.Println(err.Error())
		return
	}

	var sendbuf [comm.MSG_HEAD_LEN]byte

	sendbuf[0] = 255
	sendbuf[comm.MSG_HEAD_LEN-1] = 254

	for i := 0; i < 100000; i++ {
		err = c.SendMsg(0, sendbuf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		clientsendcnt++
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
