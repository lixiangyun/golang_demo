package srpc

import (
	"sync"
)

type Server struct {
	Addr     string
	function map[string]func()

	lock sync.Mutex
	wait sync.WaitGroup
}

func NewServer(addr string) *Server {
	return &Server{Addr: addr}
}

func (s *Server) AddMethod(fun func()) {

	s.lock.Lock()
	defer s.lock.Unlock()

	name := string(fun)
	_, b := s.function[name]
	if b == false {
		s.function[name] = fun
	}
}

func (s *Server) DelMethod(fun func()) {

	s.lock.Lock()
	defer s.lock.Unlock()

	name := string(fun)
	_, b := s.function[name]
	if b == true {
		delete(s.function, name)
	}
}

func (s *Server) Start() {

}
