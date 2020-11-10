package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
)

func GetUUID() string {
	return uuid.New().String()
}

type Clients struct {
	Address string `json:"address"`
	UUID    string `json:"uuid"`
}

type AddressList struct {
	List []Clients `json:"list"`
}

var gClients AddressList

func init()  {
	gClients.List = make([]Clients, 0)
}

func AddClient(id string, address string)  {
	for idx, v:= range gClients.List {
		if v.UUID == id {
			gClients.List[idx].Address = address
			return
		}
	}
	gClients.List = append(gClients.List, Clients{UUID: id, Address: address})
}

func SyncClient(body []byte)  {
	err := json.Unmarshal(body, &gClients)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetClient() []byte {
	body, err := json.Marshal(&gClients)
	if err != nil {
		fmt.Println(err.Error())
	}
	return body
}

func Server(port string) {
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn, err2 := net.ListenUDP("udp", addr)
	if err2 != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	var buf [65535]byte
	for {
		cnt, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			if err == io.EOF {
				fmt.Println("close connect! ", conn.RemoteAddr())
				return
			}
		}
		fmt.Printf("client : %s, address: %s\n", string(buf[:cnt]), addr.String())
		AddClient(string(buf[:cnt]), addr.String())
		conn.WriteToUDP(GetClient(), addr)
	}
}

var clientID string

func Client2Client(conn *net.UDPConn, address string, body []byte)  {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = conn.WriteToUDP(body, addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("send to client %s success\n", addr.String())
}

func Client(addr string, local string) {
	clientID = GetUUID()

	localaddr, err := net.ResolveUDPAddr("udp", local)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	remoteaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn, err := net.ListenUDP("udp", localaddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	fmt.Println( "server : " + remoteaddr.String())

	go func() {
		for {
			writebuff := []byte(clientID)
			cnt, err := conn.WriteToUDP(writebuff[:], remoteaddr)
			if err != nil {
				fmt.Println(err.Error())
			}
			if cnt > 0 {
				fmt.Printf("send to server %s success\n", remoteaddr.String())
			}
			time.Sleep(5*time.Second)
		}
	}()

	go func() {
		for  {
			for _, v:= range gClients.List {
				if v.UUID != clientID {
					body := fmt.Sprintf("hello world from: %s", clientID)
					Client2Client(conn, v.Address, []byte(body))
				}
			}
			time.Sleep(5*time.Second)
		}
	}()

	var readbuff [1024]byte

	for  {
		cnt, addr, err := conn.ReadFromUDP(readbuff[:])
		if addr.String() == remoteaddr.String() {
			SyncClient(readbuff[:cnt])
		} else {
			fmt.Printf("client recv : %s %s\n", addr.String(), string(readbuff[:cnt]))
		}
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage: <-s/-c> <ip:port>")
		return
	}
	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2],args[3])
	}
}
