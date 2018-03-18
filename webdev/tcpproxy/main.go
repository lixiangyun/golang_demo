package main

import (
	"log"
)

func main() {
	proxy := NewTcpProxy(":20022", "192.168.0.10:22")

	err := proxy.Start()
	if err != nil {
		log.Println(err.Error())
	}
}
