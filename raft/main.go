package main

import (
	"golang_demo/raft/src"
	"log"
	"os"
	"time"
)

func main() {

	args := os.Args
	if len(args) < 4 {
		log.Println("Usage: <IP:PORT> <IP:PORT> <IP:PORT> ...")
		return
	}

	r, err := raft.NewRaft(args[1], args[2:])
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = raft.Start(r)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Server start ok!")

	for {
		// run forever not to stop
		time.Sleep(time.Duration(10 * time.Second))
	}

	log.Println("Server stop ok!")

	raft.Stop(r)
}
