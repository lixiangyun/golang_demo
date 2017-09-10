package main

import (
	"fmt"
	"golang_demo/srpc/srpc"
	"os"
	"time"
)

type SAVE struct {
}

func (s *SAVE) Add(a uint32, b *uint32) error {

	*b = (a + 1)

	fmt.Println("call add ", a, b)

	return nil
}

func (s *SAVE) Sub(a uint32, b *uint32) error {

	*b = (a - *b)

	fmt.Println("call sub ", a, b)

	return nil
}

func Server(addr string) {

	var s SAVE
	server := srpc.NewServer(":1234")
	server.BindMethod(&s)

	err := server.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		time.Sleep(time.Second)
	}

	server.Stop()
}

func Client(addr string) {

	client := srpc.NewClient("localhost:1234")

	var a, b uint32
	a = 1
	b = 2
	err := client.Call("Add", a, &b)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("a=", a, " b=", b)
	}
}

func main() {

	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: < -s PORT / -c IP:PORT >")
		return
	}

	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2])
	}
}