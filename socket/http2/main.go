package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Println("Usage: <-s/-c>")
		return
	}

	switch args[1] {
	case "-s11":
		HttpServer1()
	case "-s20":
		HttpServer2()
	case "-c11":
		HttpClient1("https://localhost:9001/v1")
		HttpClient1("https://localhost:9000/v1")
		HttpClient1("https://localhost:9001/")
		HttpClient1("https://localhost:9000/")
	case "-c20":
		HttpClient2("https://localhost:9001/v1")
		HttpClient2("https://localhost:9000/v1")
		HttpClient2("https://localhost:9001/")
		HttpClient2("https://localhost:9000/")
	}
}
