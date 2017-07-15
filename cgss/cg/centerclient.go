package cg

import (
	"encoding/json"
	"errors"

	"golang_demo/cgss/ipc"
)

type CenterClient struct {
	*ipc.IpcClient
}

func (c *CenterClient) AddPlayer(p *Player) error {
	b, err := json.Marshal(*p)

	if err != nil {
		return err
	}

	resp, err := c.Call("addplayer", string(b))
	if err == nil && resp.Code == "200" {
		return nil
	}

	return err
}

func (c *CenterClient) RemovePlayer(name string) error {
	ret, _ := c.Call("removeplayer", name)

	if ret.Code == "200" {
		return nil
	}

	return errors.New(ret.Code)
}

func (c *CenterClient) ListPlayer(p string) (ps []*Player, err error) {
	resp, _ := c.Call("listplayer", p)
	if resp.Code != "200" {
		err = errors.New(resp.Code)
		return
	}

	err = json.Unmarshal([]byte(resp.Body), &ps)
	return
}

func (c *CenterClient) Broadcast(msg string) error {
	m := &Message{Content: msg} // 构造message结构体内容

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, _ := c.Call("broadcast", string(b))
	if resp.Code == "200" {
		return nil
	}

	return errors.New(resp.Code)
}
