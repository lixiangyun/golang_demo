package main

import (
	"log"
	"net"
	"sync"
)

type TcpBridge struct {
	RemoteAddr string
	LocalAddr  string
}

func NewTcpBridge(localaddr string, remoteaddr string) *TcpBridge {
	return &TcpBridge{LocalAddr: localaddr, RemoteAddr: remoteaddr}
}

func (t *TcpBridge) Start() error {

	localconn, err := net.Dial("tcp", t.LocalAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	remoteconn, err := net.Dial("tcp", t.RemoteAddr)
	if err != nil {
		log.Println(err.Error())
		localconn.Close()
		return err
	}

	syncSem := new(sync.WaitGroup)

	log.Println("new bridge connect. ", localconn.RemoteAddr(), " <-> ", remoteconn.RemoteAddr())

	syncSem.Add(2)

	go tcpChannel(localconn, remoteconn, syncSem)
	go tcpChannel(remoteconn, localconn, syncSem)

	syncSem.Wait()

	log.Println("close bridge connect. ", localconn.RemoteAddr(), " <-> ", remoteconn.RemoteAddr())

	return nil
}
