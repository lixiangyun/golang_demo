package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"
	"time"

	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

var stat_cnt, stat_size int

func init()  {
	go func() {
		for  {
			time.Sleep(time.Second)
			log.Printf("%d , %d MB/s",
				stat_cnt, stat_size/(1024*1024))
			stat_cnt = 0
			stat_size = 0
		}
	}()
}

func ReadAndWrite(conn net.Conn)  {
	defer func() {
		log.Println("connect %s close", conn.RemoteAddr().String())
		conn.Close()
	}()

	var buff [65535]byte
	for  {
		cnt, err := conn.Read(buff[:])
		if err != nil {
			return
		}
		stat_cnt++
		stat_size += cnt

		_, err = conn.Write(buff[:cnt])
		if err != nil {
			return
		}
	}
}

func server() {
	// Prepare the IP to connect to
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4444}

	// Generate a certificate and private key to secure the connection
	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Create parent context to cleanup handshaking connections on exit.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, 30*time.Second)
		},
	}

	// Connect to a DTLS server
	listener, err := dtls.Listen("udp", addr, config)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer listener.Close()

	log.Println("Listening")

	for  {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			return
		}
		go ReadAndWrite(conn)
	}
}

func server_notls() {
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4444}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()
	var buff [65535]byte

	for  {
		cnt, remote, err := conn.ReadFromUDP(buff[:])
		if err != nil {
			log.Println(err.Error())
			continue
		}
		stat_cnt++
		stat_size += cnt

		_, err = conn.WriteToUDP(buff[:cnt], remote)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func client_notls() {
	dtlsConn, err := net.Dial("udp", "127.0.0.1:4444")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer dtlsConn.Close()
	dtlsConn.Write(make([]byte, 3200))
	ReadAndWrite(dtlsConn)
}

func client() {
	// Prepare the IP to connect to
	addr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4444}

	// Generate a certificate and private key to secure the connection
	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	// Connect to a DTLS server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dtlsConn, err := dtls.DialWithContext(ctx, "udp", addr, config)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer dtlsConn.Close()

	dtlsConn.Write(make([]byte, 3200))

	ReadAndWrite(dtlsConn)
}

func main()  {
	if os.Args[1] == "-c" {
		client_notls()
	} else {
		server_notls()
	}
}