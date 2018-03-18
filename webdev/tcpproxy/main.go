package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// 实际中应该用更好的变量名
var (
	mode string

	remoteaddr string
	localaddr  string
	help       bool
)

func usage() {
	fmt.Fprintf(os.Stderr, `tcpproxy version: tcpproxy/0.1.0
Usage: tcpproxy [-hm] [-local ip:port] [-remote ip:port]

Options:
`)
	flag.PrintDefaults()
}

func init() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&mode, "m", "proxy", "using bridge/link/proxy mode.")

	flag.StringVar(&remoteaddr, "remote", "", "remote addr")
	flag.StringVar(&localaddr, "local", "", "local addr")

	flag.Usage = usage
}

func main() {

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if mode == "proxy" {

		if remoteaddr != "" && localaddr != "" {
			proxy := NewTcpProxy(localaddr, remoteaddr)

			err := proxy.Start()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

	} else if mode == "bridge" {

		if remoteaddr != "" && localaddr != "" {
			proxy := NewTcpBridge(localaddr, remoteaddr)

			err := proxy.Bridge()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
	} else if mode == "link" {

		if remoteaddr != "" && localaddr != "" {
			proxy := NewTcpBridge(localaddr, remoteaddr)

			err := proxy.Link()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
	}

	flag.Usage()
}
