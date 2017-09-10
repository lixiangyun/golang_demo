package srpc

import (
	"errors"
	"log"
	"net"
	"reflect"
	"runtime/debug"
	"time"
)

func Log(v ...interface{}) {
	log.Println(v)
	log.Println(string(debug.Stack()))
}

type Client struct {
	Addr  string
	MsgId uint64
}

func NewClient(addr string) *Client {
	var client = Client{Addr: addr}
	client.MsgId = uint64(time.Now().Nanosecond())
	return &client
}

func (n *Client) Call(method string, req interface{}, rsp interface{}) error {

	// 创建udp协议的socket服务
	socket, err := net.Dial("udp", n.Addr)
	if err != nil {
		return err
	}

	defer socket.Close()

	n.MsgId++

	var reqblock RequestBlock
	var rspblock RsponseBlock

	reqblock.MsgType = 0
	reqblock.Method = method
	reqblock.MsgId = n.MsgId
	reqblock.Parms[0] = reflect.ValueOf(req).String()
	reqblock.Parms[1] = reflect.ValueOf(rsp).String()
	reqblock.Body, err = CodePacket(req)
	if err != nil {
		Log(err.Error())
		return err
	}

	// 设置 read/write 超时时间
	err = socket.SetDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		Log(err.Error())
		return err
	}

	// 序列化请求报文
	newbuf, err := CodePacket(reqblock)
	if err != nil {
		Log(err.Error())
		return err
	}

	// 发送到服务端
	_, err = socket.Write(newbuf)
	if err != nil {
		Log(err.Error())
		return err
	}

	var buf [4096]byte

	// 获取服务端应答报文
	cnt, err := socket.Read(buf[0:])
	if err != nil {
		Log(err.Error())
		return err
	}

	// 反序列化报文
	err = DecodePacket(buf[0:cnt], rspblock)
	if err != nil {
		Log(err.Error())
		return err
	}

	// 校验请求的序号是否一致
	if rspblock.MsgId != reqblock.MsgId {
		err = errors.New("recv a bad packet ")
		Log(err.Error())
		return err
	}

	// 反序列化报文
	err = DecodePacket(rspblock.Body, rsp)
	if err != nil {
		Log(err.Error())
		return err
	}

	return nil
}
