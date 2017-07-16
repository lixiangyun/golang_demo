package ipc

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

type Response struct {
	Code string `json:"code"`
	Body string `json:"body"`
}

type Server interface {
	Name() string
	Handle(m, p string) *Response
}

type IpcServer struct {
	Server
}

func NewIpcServer(s Server) *IpcServer {
	return &IpcServer{s}
}

func (s *IpcServer) Response(c chan string) {
	for {
		request := <-c
		if request == "CLOSE" {
			fmt.Println("IPC close")
			break
		}

		var req Request

		err := json.Unmarshal([]byte(request), &req)
		if err != nil {
			fmt.Println("Invalid request format:", request)
			return
		}

		fmt.Println("request: ", req)

		resp := s.Server.Handle(req.Method, req.Params)

		fmt.Println("response: ", resp)

		b, err := json.Marshal(resp)

		c <- string(b)
	}

	fmt.Println("Session closed.")
}

func (s *IpcServer) Connect() chan string {
	session := make(chan string, 0)

	go s.Response(session)

	fmt.Println("A new session has been created sucessfully.")

	return session
}
