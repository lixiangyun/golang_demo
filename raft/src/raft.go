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
	RPC_TIMEOUT  = 500 * time.Millisecond
	VOTE_TIMEOUT = 3 * time.Second
	KEEP_ALIVE   = 1 * time.Second
	PING_TIMEOUT = 2 * time.Second
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
	votecnt   int

	client map[string]*rpc.Client

	currentTerm uint64

	ping    *time.Timer
	timeout *time.Timer
	mlock   *sync.Mutex

	lis net.Listener

	wait *sync.WaitGroup
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

	return VOTE_TIMEOUT + time.Duration(val)*time.Millisecond
}

func (r *RAFT) RequestVote(l *LeaderInfo, b *bool) error {

	log.Println("Recv RequestVote from Leader :", l.ServerName)

	r.mlock.Lock()
	defer r.mlock.Unlock()

	*b = false

	switch r.role {
	case ROLE_CANDIDATE:
		{
			if l.Term >= r.currentTerm && r.leader == "" {

				r.leader = l.ServerName
				r.role = ROLE_FOLLOWER
				r.currentTerm = l.Term

				r.timeout.Reset(getDymicTimeOut())
				log.Println("change role to follower from Leader :", l.ServerName)

				*b = true
			} else {
				log.Printf("RequestVote Name : %s , Role: %s , Term %d , Leader %s \r\n", r.name, r.role, r.currentTerm, r.leader)
			}
		}
	case ROLE_FOLLOWER:
		{
			if l.Term >= r.currentTerm && r.leader == "" {

				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.timeout.Reset(getDymicTimeOut())

				log.Println("agree new Leader :", l.ServerName)
				*b = true
			} else {
				log.Printf("RequestVote Name : %s , Role: %s , Term %d , Leader %s \r\n", r.name, r.role, r.currentTerm, r.leader)
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
		} else {
			log.Printf("AppendEntries Name : %s , Role: %s , Term %d , Leader %s \r\n", r.name, r.role, r.currentTerm, r.leader)
		}
	} else {
		log.Printf("AppendEntries Name : %s , Role: %s , Term %d , Leader %s \r\n", r.name, r.role, r.currentTerm, r.leader)
	}

	return nil
}

func (r *RAFT) Heartbeats(l *LeaderInfo, b *bool) error {

	log.Println("Recv Heartbeats from Leader :", l.ServerName)

	*b = false

	r.mlock.Lock()
	defer r.mlock.Unlock()

	switch r.role {
	case ROLE_FOLLOWER:
		{
			if l.ServerName == r.leader {

				r.timeout.Reset(getDymicTimeOut())
				log.Println("Get Heartbeats from Leader :", l.ServerName)

				*b = true
			}
		}
	case ROLE_CANDIDATE:
		{
			if l.Term >= r.currentTerm {

				r.role = ROLE_FOLLOWER
				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.timeout.Reset(getDymicTimeOut())

				log.Println("follower new Leader :", l.ServerName)
			}
		}
	case ROLE_LEADER:
		{
			if l.Term > r.currentTerm {

				r.role = ROLE_FOLLOWER
				r.leader = l.ServerName
				r.currentTerm = l.Term
				r.timeout.Reset(getDymicTimeOut())

				log.Println("follower new Leader :", l.ServerName)
			}
		}
	}

	return nil
}

func Ping(r *RAFT) {
	log.Println("Ping Start")

	defer r.wait.Done()
	defer log.Println("Ping Close")

	for {
		_, b := <-r.ping.C
		if b == false {
			return
		}

		for _, name := range r.othername {
			_, b := r.client[name]
			if b == false {
				client, err := rpc.DialHTTP("tcp", name)
				if err != nil {
					log.Println("Ping ", err.Error())
					continue
				}
				r.client[name] = client
			}
		}
	}
}

func TimeOut(r *RAFT) {

	log.Println("Timer Start ")

	defer r.wait.Done()
	defer log.Println("Timer Close")

	for {
		_, b := <-r.timeout.C
		if b == false {
			return
		}

		log.Printf("Now Self : %s , Role: %s , Term %d , Leader %s \r\n", r.name, r.role, r.currentTerm, r.leader)

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
				r.mlock.Unlock()

				LeaderVote(r)
			}
		default:
			log.Println("Recv Stop Raft")
			return
		}
	}
}

func keepalive(r *RAFT, name string, cli *rpc.Client, tm time.Duration, b chan bool) {
	rpctm := time.NewTimer(RPC_TIMEOUT)

	request := new(LeaderInfo)
	response := new(bool)

	request.ServerName = r.name
	request.Term = r.currentTerm

	result := cli.Go("RAFT.Heartbeats", request, response, nil)

	select {
	case rsp := <-result.Done:
		{
			if rsp.Error != nil {
				cli.Close()
				r.client[]
			}
		}
	case <-rpctm:
		{

		}
	}
}

func KeepAlive(r *RAFT) {

	log.Println("KeepAlive timer start ")

	length := len(r.othername)
	queue := make(chan *rpc.Call, length)

	for _, cli := range r.client {
		if cli == nil {
			continue
		}

	}

	agreenum := 0

	for i := 0; i < length; i++ {

		select {
		case result := <-queue:
			{

			}
		case <-rpctm:
			{

			}
		}

		result := <-queue

	}

	close(queue)

	if agreenum == 0 {
		r.mlock.Lock()

		r.leader = ""
		r.role = ROLE_FOLLOWER
		r.timeout.Reset(getDymicTimeOut())

		r.mlock.Unlock()
	} else {

		r.mlock.Lock()
		r.timeout.Reset(KEEP_ALIVE)
		r.mlock.Unlock()
	}
}

func LeaderVote(r *RAFT) {
	var err error

	log.Println("Start LeaderVote for ", r.name)

	totalnum := len(r.othername)
	agreenum := 0

	r.mlock.Lock()
	if r.leader != "" {
		r.mlock.Unlock()

		log.Println("Stop LeaderVote for ", r.leader)
		return
	}
	r.leader = r.name
	r.mlock.Unlock()

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

		err = cli.Call("RAFT.RequestVote", &request, &response)
		if err != nil {
			log.Println("LeaderVote ", err.Error())
			continue
		}
	}

	for _, cli := range client {

		if cli == nil {
			continue
		}

		request.ServerName = r.name
		request.Term = r.currentTerm

		err = cli.Call("RAFT.AppendEntries", &request, &response)
		if err != nil {
			log.Println("LeaderVote ", err.Error())
			continue
		}

		if response == true {
			agreenum++
		} else {
			log.Println("LeaderVote not agree!")
		}
	}

	if agreenum*2 >= totalnum {

		r.mlock.Lock()

		r.role = ROLE_LEADER
		r.timeout.Reset(KEEP_ALIVE)

		r.mlock.Unlock()

		log.Println("LeaderVote success! ")
	} else {
		r.mlock.Lock()

		r.leader = ""
		r.timeout.Reset(getDymicTimeOut())

		r.mlock.Unlock()

		log.Println("LeaderVote failed! ")
	}

	for _, cli := range client {
		if cli == nil {
			continue
		}
		cli.Close()
	}

	log.Println("LeaderVote timer close ")
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
	r.votecnt = 0

	r.wait = new(sync.WaitGroup)

	r.mlock = new(sync.Mutex)

	return r, nil
}

func Start(r *RAFT) error {

	log.Println("Server Start")

	server := rpc.NewServer()

	err := server.Register(r)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	listen, err := net.Listen("tcp", ":"+r.port)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	r.lis = listen

	r.wait.Add(3)

	r.ping = time.NewTimer(PING_TIMEOUT)
	go Ping(r)

	r.timeout = time.NewTimer(getDymicTimeOut())
	go TimeOut(r)

	go func() {

		defer r.wait.Done()

		err = http.Serve(listen, nil)

		log.Println("Server End")

		if err != nil {
			log.Println(err.Error())
			return
		}
	}()

	return nil
}

func Stop(r *RAFT) {

	r.mlock.Lock()
	r.role = ""
	r.lis.Close()
	r.timeout.Reset(0)
	r.mlock.Unlock()

	r.wait.Wait()
}
