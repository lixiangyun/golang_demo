package raft

import (
	"crypto/rand"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

const (
	ROLE_FOLLOWER  = "Follower"
	ROLE_CANDIDATE = "Candidate"
	ROLE_LEADER    = "Leader"
)

const (
	TIMEOUT   = 3 * time.Second
	KEEPALIVE = 1 * time.Second
)

type LeaderInfo struct {
	ServerName string
	Term       uint64
}

type RAFT struct {
	role      string
	port      string
	name      string
	othername []string
	leader    string

	currentTerm uint64

	timeout *time.Timer
	mlock   *sync.Mutex
}

func getDymicTimeOut() time.Duration {

	var buf [4]byte
	var val int

	_, err := rand.Read(buf[0:])
	if err != nil {
		val = time.Now().Nanosecond()
	} else {
		val = int(buf[0]) + int(buf[1])<<8 + int(buf[2])<<16 + int(buf[3])<<24
	}

	val = val % 3000

	log.Println("val : ", val)

	return TIMEOUT + time.Duration(val)*time.Millisecond
}

func (r *RAFT) RequestVote(l *LeaderInfo, b *bool) error {

	log.Println("Recv RequestVote from Leader :", l.ServerName)

	r.mlock.Lock()
	defer r.mlock.Unlock()

	*b = false

	switch r.role {
	case ROLE_CANDIDATE:
		{
			if l.Term > r.currentTerm && r.leader == "" {

				r.leader = l.ServerName
				r.role = ROLE_FOLLOWER
				r.currentTerm = l.Term

				r.timeout.Reset(getDymicTimeOut())
				log.Println("change role to follower from Leader :", l.ServerName)

				*b = true
			}
		}
	case ROLE_FOLLOWER:
		{
			if l.Term > r.currentTerm && r.leader == "" {

				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.timeout.Reset(getDymicTimeOut())

				log.Println("agree new Leader :", l.ServerName)
				*b = true
			}
		}
	}
	return nil
}

func (r *RAFT) AppendEntries(l *LeaderInfo, b *bool) error {

	log.Println("Recv AppendEntries from Leader :", l.ServerName)

	r.mlock.Lock()
	defer r.mlock.Unlock()

	*b = false
	if r.role == ROLE_FOLLOWER {
		if r.leader == l.ServerName {

			r.timeout.Reset(getDymicTimeOut())

			log.Println("agree append entries request from Leader :", l.ServerName)

			*b = true
		}
	}

	return nil
}

func (r *RAFT) Heartbeats(h *LeaderInfo, b *bool) error {

	log.Println("Recv Heartbeats from Leader :", h.ServerName)

	*b = false

	r.mlock.Lock()
	defer r.mlock.Unlock()

	if r.role == ROLE_FOLLOWER {
		if h.ServerName == r.leader {
			r.timeout.Reset(getDymicTimeOut())
			log.Println("Get Heartbeats from Leader :", h.ServerName)

			*b = true
		}
	}

	return nil
}

func TimeOut(r *RAFT) {

	log.Println("Timer Start ")

	for {
		_, b := <-r.timeout.C
		if b == false {
			log.Println("Timer Close.")
			return
		}

		log.Printf("Name : %s , Role: %s , Term %d \r\n", r.name, r.role, r.currentTerm)

		switch r.role {
		case ROLE_LEADER:
			{
				KeepAlive(r)
			}
		case ROLE_CANDIDATE:
			{
				LeaderVote(r)
			}
		case ROLE_FOLLOWER:
			{
				r.mlock.Lock()

				r.role = ROLE_CANDIDATE
				r.leader = ""
				r.currentTerm++
				r.timeout.Reset(getDymicTimeOut())

				r.mlock.Unlock()
			}
		default:
			log.Println("unknow role", r.role)
		}
	}

	log.Println("Timer Close.")
}

func KeepAlive(r *RAFT) {

	log.Println("KeepAlive timer start ")

	agreenum := 0
	for _, v := range r.othername {

		client, err := rpc.DialHTTP("tcp", v)
		if err != nil {
			log.Println("KeepAlive ", err.Error())
			continue
		}

		defer client.Close()

		var request LeaderInfo
		var response bool

		request.ServerName = r.name
		request.Term = r.currentTerm

		err = client.Call("RAFT.Heartbeats", &request, &response)
		if err != nil {
			log.Println("KeepAlive ", err.Error())
			continue
		}

		if response == true {
			agreenum++
			continue
		}
	}

	if agreenum == 0 {
		r.mlock.Lock()

		r.leader = ""
		r.role = ROLE_FOLLOWER
		r.timeout.Reset(getDymicTimeOut())

		r.mlock.Unlock()
	} else {

		r.mlock.Lock()
		r.timeout.Reset(KEEPALIVE)
		r.mlock.Unlock()
	}
}

func LeaderVote(r *RAFT) {
	var err error

	totalnum := len(r.othername)
	agreenum := 0

	client := make([]*rpc.Client, totalnum)

	for i, v := range r.othername {

		client[i], err = rpc.DialHTTP("tcp", v)
		if err != nil {
			log.Println("LeaderVote ", err.Error())
			continue
		}
	}

	var request LeaderInfo
	var response bool

	for _, cli := range client {

		if cli == nil {
			continue
		}

		request.ServerName = r.name
		request.Term = r.currentTerm

		//err = client.Go()

		err = cli.Call("RAFT.RequestVote", &request, &response)
		if err != nil {
			log.Println(" LeaderVote ", err.Error())
			continue
		}
	}

	for _, cli := range client {

		if cli == nil {
			continue
		}

		request.ServerName = r.name
		request.Term = r.currentTerm

		//err = client.Go()

		err = cli.Call("RAFT.AppendEntries", &request, &response)
		if err != nil {
			log.Println(" LeaderVote ", err.Error())
			continue
		}

		if response == true {
			agreenum++
		} else {
			log.Println(" LeaderVote not agree!")
		}
	}

	if agreenum*2 >= totalnum {

		r.mlock.Lock()

		r.leader = r.name
		r.role = ROLE_LEADER
		r.timeout.Reset(KEEPALIVE)

		r.mlock.Unlock()

		log.Println(" LeaderVote success! ")
	} else {
		r.mlock.Lock()
		r.timeout.Reset(getDymicTimeOut())
		r.mlock.Unlock()
	}

	for _, cli := range client {
		if cli == nil {
			continue
		}
		cli.Close()
	}

	log.Println(" LeaderVote timer close ")
}

func NewRaft(selfaddr string, otheraddr []string) (*RAFT, error) {

	idx := strings.Index(selfaddr, ":")
	if idx == -1 {
		return nil, errors.New("Input Selfaddr invailed : " + selfaddr)
	}

	r := new(RAFT)

	r.currentTerm = 0
	r.leader = ""
	r.name = selfaddr
	r.othername = otheraddr
	r.port = selfaddr[idx+1:]
	r.role = ROLE_FOLLOWER

	r.mlock = new(sync.Mutex)

	return r, nil
}

func Start(r *RAFT) {

	log.Println("Server Start")

	server := rpc.NewServer()

	err := server.Register(r)
	if err != nil {
		log.Println(err.Error())
		return
	}

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	listen, err := net.Listen("tcp", ":"+r.port)
	if err != nil {
		log.Println(err.Error())
		return
	}

	r.timeout = time.NewTimer(getDymicTimeOut())
	go TimeOut(r)

	err = http.Serve(listen, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("Server End")
}

func Stop(r *RAFT) {

}
