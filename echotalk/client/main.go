package main

import (
	"crypto/tls"
	"io"
	"log"
)

func main() {

	config := tls.Config{InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", "localhost:9090", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
		return
	}

	defer conn.Close()

	log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	message := "Hello World\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalln("client: write: %s", err)
		return
	}

	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 512)
	n, err = conn.Read(reply)
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting..")

}
