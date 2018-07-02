package main

import (
	"log"
	"net"
)

func getLocalIp() []string {

	IpAddr := make([]string, 0)

	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		log.Println("Get local IP addr failed!!!")
		return IpAddr
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				IpAddr = append(IpAddr, ipnet.IP.String())
			}
		}
	}

	return IpAddr
}

func main() {

	ipaddr := getLocalIp()

	for _, v := range ipaddr {
		log.Println(v)
	}
}
