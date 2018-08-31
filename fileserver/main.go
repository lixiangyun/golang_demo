package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	listener string
	path     string
	tls      bool
	cert     string
	key      string
	help     bool
)

func init() {
	flag.BoolVar(&help, "h", false, "this help.")
	flag.StringVar(&listener, "p", "127.0.0.1:8080", "http file server listener address.")
	flag.StringVar(&path, "d", "./", "http file server root directory.")
	flag.BoolVar(&tls, "tls", false, "enable https.")
	flag.StringVar(&cert, "cert", "server.crt", "certificate file (https).")
	flag.StringVar(&key, "key", "server.key", "private key file name (https).")
}

func main() {

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	fileinfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !fileinfo.IsDir() {
		fmt.Println("path is not directory.", path)
		return
	}

	h := http.FileServer(http.Dir(path))

	if tls {
		log.Fatal(http.ListenAndServeTLS(listener, cert, key, h))
	} else {
		log.Fatal(http.ListenAndServe(listener, h))
	}
}
