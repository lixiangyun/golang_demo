package main

import (
	"fmt"
	"golang_demo/raft/src"
	"os"
)

func main() {

	args := os.Args
	if len(args) < 4 {
		fmt.Println("Usage: <IP:PORT> <IP:PORT> <IP:PORT> ...")
		return
	}

	r, err := raft.NewRaft(args[1], args[2:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	raft.Start(r)
}
