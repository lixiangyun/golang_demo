package cg

import (
	"encoding/json"
	"errors"
	"golang_demo/cgss/ipc"
	"sync"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

/* 这里  `...` 中的文本叫做tag，它只是普通字符串，你也可以理解为对结构体字段的注释；
   很多编解码器都会利用反射机制的这个能力做编解码(Marshal/Unmarshal), 这里是JSON
   编解码器叫做Marshal/Unmarshal需要。
*/

type ConterServer struct {
	servers map[string]ipc.Server
	players []*Player
	mutex   sync.RWMutex
}

var _ ipc.Server = &ConterServer{}

func (s *ConterServer) addPlayer(p string) error {
	player := NewPlayer()

	err := json.Unmarshal([]byte(p), &player)
	if err != nil {
		return err
	}

	s.mutex.Lock()

	defer s.mutex.Unlock()

	s.players = append(s.players, player)

	return nil
}

func (s *ConterServer) removePlayer(p string) error {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	for i, v := range s.players {
		if v.Name == p {
			s.players = append(s.players[:i-1], s.players[:i+1]...)
			return nil
		}
	}

	return errors.New("Play not found.")
}

func (s *ConterServer) listPlayer(p string) (player string, err error) {
	s.mutex.RLock()

	defer s.mutex.RUnlock()

	if len(s.players) > 0 {
		b, _ := json.Marshal(s.players)
		player = string(b)
	} else {
		err = errors.New("No player online.")
	}
	return
}

func (s *ConterServer) broadcast(p string) error {
	var msg Message

	err := json.Unmarshal([]byte(p), &msg)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.players) > 0 {
		for _, player := range s.players {
			player.mq <- &msg
		}
	} else {
		err = errors.New("No player online.")
	}
	return err
}

func (s *ConterServer) Handle(m, p string) *ipc.Response {
	switch m {
	case "addplayer":
		err := s.addPlayer(p)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
	case "removeplayer":
		err := s.removePlayer(p)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
	case "listplayer":
		players, err := s.listPlayer(p)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{"200", players}
	case "broadcast":
		err := s.broadcast(p)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{Code: "200"}
	default:
		return &ipc.Response{Code: "404", Body: m + ":" + p}
	}
	return &ipc.Response{Code: "200"}
}

func (s *ConterServer) Name() string {
	return "CenterServer"
}
