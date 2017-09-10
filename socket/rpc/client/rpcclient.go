package main

import (
	"fmt"
	"net/rpc"

	"golang_demo/socket/rpc/server"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	args := &server.Args{7, 8}

	var reply int

	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Arith.Multiply: %d*%d=%d \r\n", args.A, args.B, reply)

	quo := new(server.Quotient)
	divcall := client.Go("Arith.Divide", args, &quo, nil)
	replycall := <-divcall.Done

	if replycall.Error != nil {
		fmt.Println(replycall.Error.Error())
		return
	}

	fmt.Printf("Arith.Divide: %d/%d=%d %d \r\n", args.A, args.B, quo.Quo, quo.Rem)

}
