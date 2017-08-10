package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"time"
)

const (
	ROLE_FOLLOWER  = 0
	ROLE_CANDIDATE = 1
	ROLE_LEADER    = 2
)

const (
	TIMEOUT   = 5 * time.Second
	KEEPALIVE = 1 * time.Second
	LEADVOTE  = 127 * time.Millisecond
)

type LeaderInfo struct {
	ServerName string
	Term       uint64
}

type HeartInfo struct {
	Servername string
	timeout    time.Duration
}

type RAFT struct {
	role      int
	port      string
	name      string
	othername map[string]string
	leader    string

	currentTerm uint64

	timeout    *time.Timer
	keepAlive  *time.Timer
	leaderVote *time.Timer
}

var raft RAFT

func (r *RAFT) RequestVote(l *LeaderInfo, b *bool) error {
	switch r.role {
	case ROLE_LEADER:
		{
			*b = false
		}

	case ROLE_CANDIDATE:
		{
			if l.Term > r.currentTerm && r.leader == "" {

				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.role = ROLE_FOLLOWER
				r.timeout = time.NewTimer(TIMEOUT)
				go TimeOut(raft.timeout)

				r.leaderVote.Stop()

				*b = true

			} else {
				*b = false
			}
		}

	case ROLE_FOLLOWER:
		{
			if r.leader == "" {
				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.timeout.Reset(TIMEOUT)

				*b = true
			} else {
				*b = false
			}
		}
	default:
		*b = false
	}
	return nil
}

func (r *RAFT) AppendEntries(l *LeaderInfo, b *bool) error {
	switch r.role {
	case ROLE_LEADER:
		{
			*b = false
		}

	case ROLE_CANDIDATE:
		{
			if l.Term > r.currentTerm && r.leader == "" {

				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.role = ROLE_FOLLOWER
				r.timeout = time.NewTimer(TIMEOUT)
				go TimeOut(raft.timeout)

				r.leaderVote.Stop()

				*b = true

			} else {
				*b = false
			}
		}

	case ROLE_FOLLOWER:
		{
			if r.leader == "" {
				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.timeout.Reset(TIMEOUT)

				*b = true
			} else {
				*b = false
			}
		}
	default:
		*b = false
	}

	return nil
}

func (r *RAFT) Heartbeats(h *LeaderInfo, b *bool) error {

	switch r.role {
	case ROLE_FOLLOWER:
		{
			if h.ServerName == r.leader {
				*b = true
				r.timeout.Reset(TIMEOUT)
				fmt.Println("Get Heartbeats from Leader :", h.ServerName)
			} else {
				*b = false
			}

		}
	default:
		*b = false
	}

	return nil
}

func ServerStart(port string) {

	server := rpc.NewServer()

	err := server.Register(&raft)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = http.Serve(listen, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func TimeOut(tm *time.Timer) {

	_, b := <-tm.C
	if b == false {
		fmt.Println(time.Now(), " TimeOut Timer Close.")
		return
	}

	raft.leader = ""

	tm.Stop()

	raft.leaderVote = time.NewTimer(321 * time.Millisecond)

	go LeaderVote(raft.leaderVote)

	fmt.Println(time.Now(), " TimeOut timer close ")
}

func KeepAlive(tm *time.Timer) {

	var agreenum int

	for {
		_, b := <-tm.C
		if b == false {
			fmt.Println(time.Now(), " KeepAlive Timer Close.")
			break
		}

		agreenum = 0

		for _, v := range raft.othername {
			client, err := rpc.DialHTTP("tcp", v)
			if err != nil {
				fmt.Println(time.Now(), " KeepAlive ", err.Error())
				continue
			}

			defer client.Close()

			var request LeaderInfo
			var response bool

			request.ServerName = raft.name
			request.Term = raft.currentTerm

			err = client.Call("RAFT.Heartbeats", &request, &response)
			if err != nil {
				fmt.Println(time.Now(), " KeepAlive ", err.Error())
				continue
			}

			if response == true {
				agreenum++
				continue
			}

			err = client.Call("RAFT.AppendEntries", &request, &response)
			if err != nil {
				fmt.Println(time.Now(), " KeepAlive ", err.Error())
				continue
			}
		}

		if agreenum == 0 {
			// 孤岛
			tm.Stop()

			raft.leader = ""
			raft.role = ROLE_FOLLOWER
			raft.timeout = time.NewTimer(TIMEOUT)
			go TimeOut(raft.timeout)

			break
		} else {
			tm.Reset(KEEPALIVE)
		}
	}

	fmt.Println(time.Now(), " KeepAlive timer close ")
}

func LeaderVote(tm *time.Timer) {

	var totalnum int
	var agreenum int

	totalnum = len(raft.othername)

	for {
		_, b := <-tm.C
		if b == false {
			fmt.Println(time.Now(), " LeaderVote channal close ")
			break
		}

		raft.currentTerm++

		agreenum = 0
		for _, v := range raft.othername {

			client, err := rpc.DialHTTP("tcp", v)
			if err != nil {
				fmt.Println(time.Now(), " LeaderVote ", err.Error())
				continue
			}

			defer client.Close()

			var request LeaderInfo
			var response bool

			request.ServerName = raft.name
			request.Term = raft.currentTerm

			err = client.Call("RAFT.RequestVote", &request, &response)
			if err != nil {
				fmt.Println(time.Now(), " LeaderVote ", err.Error())
				continue
			}

			if response == true {
				agreenum++
			} else {
				fmt.Println(time.Now(), " LeaderVote ", v, " not agree!")
			}
		}

		if agreenum*2 >= totalnum {

			raft.leader = raft.name
			raft.role = ROLE_LEADER

			fmt.Println(time.Now(), " LeaderVote success! ")

			tm.Stop()

			raft.keepAlive = time.NewTimer(KEEPALIVE)
			go KeepAlive(raft.keepAlive)

			break
		}
	}

	fmt.Println(time.Now(), " LeaderVote timer close ")
}

func main() {

	args := os.Args
	if len(args) < 4 {
		fmt.Println("Usage: <IP:PORT> <IP:PORT> <IP:PORT> ...")
		return
	}

	raft.othername = make(map[string]string, len(args)-1)

	for _, v := range args[2:] {
		raft.othername[v] = v
	}

	raft.name = args[1]
	raft.role = ROLE_FOLLOWER
	raft.leader = ""

	index := strings.Index(raft.name, ":")
	if index == -1 {
		fmt.Println("Input addr invailed!")
		return
	}

	raft.port = raft.name[index+1:]
	raft.currentTerm = 0

	fmt.Println(raft)

	raft.timeout = time.NewTimer(TIMEOUT)
	go TimeOut(raft.timeout)

	ServerStart(raft.port)
}
