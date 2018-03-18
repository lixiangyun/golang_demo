package main

import (
	"flag"
	"log"
)

// 实际中应该用更好的变量名
var (
	bridge     bool
	link       bool
	tcpproxy   bool
	remoteaddr string
	localaddr  string
	help       bool
)

func init() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.BoolVar(&bridge, "b", false, "using bridge mode.")
	flag.BoolVar(&link, "l", false, "using link connect mode.")
	flag.BoolVar(&tcpproxy, "t", false, "using tcp proxy mode.")

	flag.StringVar(&remoteaddr, "remote", "", "remote addr")
	flag.StringVar(&localaddr, "local", "", "local addr")
}

func main() {

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if tcpproxy {

		if remoteaddr != "" && localaddr != "" {
			proxy := NewTcpProxy(localaddr, remoteaddr)

			err := proxy.Start()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

	} else if bridge {

		if remoteaddr != "" && localaddr != "" {
			proxy := NewTcpBridge(localaddr, remoteaddr)

			err := proxy.Bridge()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
	} else if link {

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
