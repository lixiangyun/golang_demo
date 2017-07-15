package cg

import (
	"fmt"
)

type Player struct {
	Name  string
	Level int
	Exp   int
	Room  int

	mq chan *Message
}

func NewPlayer() *Player {
	m := make(chan *Message, 1024)

	p := &Player{"", 0, 0, 0, m}

	go func(p *Player) {
		for {
			msg := <-p.mq

			fmt.Println(p.Name, "recvived messages:", msg.Content)
		}
	}(p)

	return p
}
