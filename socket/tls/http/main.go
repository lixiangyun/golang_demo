package main

import (
	"log"
	"os"
)

var (
	CA_FILE     = "../crt/ca.crt"
	SERVER_CERT = "../crt/server.crt"
	SERVER_KEY  = "../crt/server.pem"
	CLIENT_CERT = "../crt/client.crt"
	CLIENT_KEY  = "../crt/client.pem"
)

func main() {

	args := os.Args

	if len(args) < 3 {
		log.Println("Usage: <-s/-c> <ip:port>")
		return
	}

	switch args[1] {
	case "-s":
		HttpServer(args[2])
	case "-c":
		HttpClient(args[2])
	}
}
