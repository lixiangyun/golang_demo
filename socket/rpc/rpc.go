package main

import (
	"fmt"

	"net"
	"net/http"
	"net/rpc"

	"golang_demo/socket/rpc/server"
)

func main() {
	arith := new(server.Arith)

	server := rpc.NewServer()

	err := server.Register(arith)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = http.Serve(lis, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
