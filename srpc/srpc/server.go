package srpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"sync"
)

type requestBlock struct {
	MsgType uint64
	MsgId   uint64
	Method  string
	Parms   [2]string
	Body    []byte
}

type rsponseBlock struct {
	MsgType uint64
	MsgId   uint64
	Method  string
	Result  error
	Body    []byte
}

type funcinfo struct {
	function reflect.Value
	functype reflect.Type
	input    [2]reflect.Type
	output   reflect.Type
}

type Server struct {
	Addr   string
	symbol map[string]funcinfo

	conn *net.UDPConn
	lock *sync.Mutex
	wait *sync.WaitGroup
}

// 报文序列化
func CodePacket(req interface{}) ([]byte, error) {
	iobuf := new(bytes.Buffer)

	err := binary.Write(iobuf, binary.BigEndian, req)
	if err != nil {
		return nil, err
	}

	//log.Println("REQ: ", req)
	//log.Println("SEND_BUF: ", iobuf.Len(), iobuf.Bytes())

	return iobuf.Bytes(), nil
}

// 报文反序列化
func DecodePacket(buf []byte, rsp interface{}) error {

	iobuf := bytes.NewReader(buf)
	err := binary.Read(iobuf, binary.BigEndian, rsp)

	//log.Println("RSP: ", rsp)
	//log.Println("RECV_BUF:", len(buf), buf)

	return err
}

func NewServer(addr string) *Server {

	s := new(Server)

	s.Addr = addr
	s.symbol = make(map[string]funcinfo, 0)

	return s
}

func (s *Server) BindMethod(pthis interface{}) {

	//创建反射变量，注意这里需要传入ruTest变量的地址；
	//不传入地址就只能反射Routers静态定义的方法
	vfun := reflect.ValueOf(pthis)
	vtype := vfun.Type()

	//读取方法数量
	num := vfun.NumMethod()

	fmt.Println("NumMethod:", num)

	//遍历路由器的方法，并将其存入控制器映射变量中
	for i := 0; i < num; i++ {

		var fun funcinfo
		fun.function = vfun.Method(i)
		fun.functype = vfun.Method(i).Type()
		funname := vtype.Method(i).Name

		if fun.functype.NumIn() != 2 {
			fmt.Printf("function %s (input parms %d) failed! \r\n", funname, fun.functype.NumIn())
			continue
		}

		if fun.functype.NumOut() != 1 {
			fmt.Printf("function %s (output parms %d) failed! \r\n", funname, fun.functype.NumOut())
			continue
		}

		fun.input[0] = fun.functype.In(0)
		fun.input[1] = fun.functype.In(1)
		fun.output = fun.functype.Out(0)

		if fun.output.String() != "error" {
			fmt.Printf("function %s (output type %s) failed! \r\n", funname, fun.output)
			continue
		}

		s.symbol[funname] = fun

		fmt.Printf("Add Method: %s \r\n", funname)
	}
}

func (s *Server) MatchMethod(method string, parms [2]string) ([]reflect.Type, error) {

	fun, b := s.symbol[method]
	if b == false {
		return nil, errors.New("can not found " + method)
	}

	for i := 0; i < 2; i++ {
		if parms[i] != fun.input[i].String() {
			errs := fmt.Sprintf("parm(%d) type not match : %s -> %s \r\n",
				i, parms[i], fun.input[i].String())
			return nil, errors.New(errs)
		}
	}

	return fun.input[0:], nil
}

func (s *Server) Call(method string, req interface{}, rsp interface{}) error {

	fun, b := s.symbol[method]
	if b == false {
		return errors.New("can not found " + method)
	}

	parms := make([]reflect.Value, 2)
	parms[0] = reflect.ValueOf(req)
	parms[1] = reflect.ValueOf(rsp)

	for i := 0; i < 2; i++ {
		if parms[i].Type() != fun.input[i] {
			errs := fmt.Sprintf("parm type not match : %s %s \r\n",
				parms[i].Type().String(), fun.input[i])
			return errors.New(errs)
		}
	}

	parms = fun.function.Call(parms)

	if len(parms) < 1 {
		return nil
	}

	if parms[0].Type().Name() == "error" {
		i := parms[0].Interface()
		if i != nil {
			return i.(error)
		}
	}

	return nil
}

// 消息收发的处理协成
func serverProcess(s *Server) {

	defer s.wait.Done()
	var buf [4096]byte

	for {
		// 监听
		n, addr, err := s.conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Println("server shutdown.")
			return
		}

		var reqblock requestBlock
		var rspblock rsponseBlock

		// 反序列化客户端请求的报文
		err = DecodePacket(buf[:n], &reqblock)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		parmtype, err := s.MatchMethod(reqblock.Method, reqblock.Parms)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		parm0 := reflect.New(parmtype[0])
		parm1 := reflect.New(parmtype[1])

		err = DecodePacket(reqblock.Body, parm0)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = s.Call(reqblock.Method, parm0.Interface(), parm1.Interface())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		rspblock.MsgType = reqblock.MsgType
		rspblock.MsgId = reqblock.MsgId
		rspblock.Method = reqblock.Method
		rspblock.Result = err

		rspblock.Body, err = CodePacket(parm1.Interface())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		rspBuf, err := CodePacket(rspblock)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 将序列化后的报文发送到客户端
		_, err = s.conn.WriteToUDP(rspBuf, addr)
		if err != nil {
			log.Println(err.Error())
			continue
		}

	}
}

func (s *Server) Start() error {

	addr, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		return err
	}

	conn, err1 := net.ListenUDP("udp", addr)
	if err1 != nil {
		return err1
	}

	s.wait.Add(1)
	s.conn = conn

	go serverProcess(s)
	return nil
}

func (s *Server) Stop() {
	s.conn.Close()
	s.wait.Wait()
}
