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

func tcpBridgeLink(localconn, remoteconn net.Conn) {

	syncSem := new(sync.WaitGroup)

	log.Println("new bridge connect. ", localconn.RemoteAddr(), " <-> ", remoteconn.RemoteAddr())

	syncSem.Add(2)

	go tcpChannel(localconn, remoteconn, syncSem)
	go tcpChannel(remoteconn, localconn, syncSem)

	syncSem.Wait()

	log.Println("close bridge connect. ", localconn.RemoteAddr(), " <-> ", remoteconn.RemoteAddr())
}

func (t *TcpBridge) Bridge() error {

	remoteconn := make(chan net.Conn, 10)
	localconn := make(chan net.Conn, 10)

	locallisten, err := net.Listen("tcp", t.LocalAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	remotelisten, err := net.Listen("tcp", t.RemoteAddr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	go func() {
		for {
			temp, err := locallisten.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			localconn <- temp
		}
	}()

	go func() {
		for {
			temp, err := remotelisten.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			remoteconn <- temp
		}
	}()

	for {
		conn1 := <-remoteconn
		conn2 := <-localconn

		go tcpBridgeLink(conn1, conn2)
	}

	return nil
}

func (t *TcpBridge) Link() error {

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
