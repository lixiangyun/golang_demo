package srpc

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Server struct {
	Addr     string
	MapsType map[string]reflect.Value

	lock *sync.Mutex
	wait *sync.WaitGroup
}

func NewServer(addr string) *Server {

	s := new(Server)

	s.Addr = addr
	s.MapsType = make(map[string]reflect.Value, 0)
	s.lock = new(sync.Mutex)
	s.wait = new(sync.WaitGroup)

	return s
}

func (s *Server) Bind(pthis interface{}) {

	s.lock.Lock()
	defer s.lock.Unlock()

	//创建反射变量，注意这里需要传入ruTest变量的地址；
	//不传入地址就只能反射Routers静态定义的方法
	vf := reflect.ValueOf(pthis)
	vft := vf.Type()

	//读取方法数量
	mNum := vf.NumMethod()
	fmt.Println("NumMethod:", mNum)
	//遍历路由器的方法，并将其存入控制器映射变量中

	for i := 0; i < mNum; i++ {

		funname := vft.Method(i).Name
		fmt.Println("index:", i, " MethodName:", funname)

		functype := vf.Method(i)
		s.MapsType[funname] = functype

		fmt.Println("input: ", functype.Type().In(0).String(), functype.Type().In(1).String())

		fmt.Println("output: ", functype.Type().Out(0).String())
	}
}

func (s *Server) Call(method string, req interface{}, rsp interface{}) error {

	v, b := s.MapsType[method]
	if b == false {
		return errors.New("can not found " + method)
	}

	parms := make([]reflect.Value, 2)
	parms[0] = reflect.ValueOf(req)
	parms[1] = reflect.ValueOf(rsp)

	parms = v.Call(parms)

	if len(parms) < 1 {
		return nil
	}

	if parms[0].Type().Name() == "error" {
		i := parms[0].Interface()
		fmt.Println("error : ", i)
	}

	return nil
}

func (s *Server) Start() {

}
