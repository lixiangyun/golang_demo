package srpc

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type funcinfo struct {
	name     string
	input    [2]string
	output   string
	function reflect.Value
	functype reflect.Type
}

type Server struct {
	Addr   string
	symbol map[string]funcinfo

	lock *sync.Mutex
	wait *sync.WaitGroup
}

func NewServer(addr string) *Server {

	s := new(Server)

	s.Addr = addr
	s.symbol = make(map[string]funcinfo, 0)

	s.lock = new(sync.Mutex)
	s.wait = new(sync.WaitGroup)

	return s
}

func (s *Server) BindMethod(pthis interface{}) {

	s.lock.Lock()
	defer s.lock.Unlock()

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
		fun.name = vtype.Method(i).Name

		if fun.functype.NumIn() != 2 {
			fmt.Printf("function %s (input parms %d) failed! \r\n", fun.name, fun.functype.NumIn())
			continue
		}

		if fun.functype.NumOut() != 1 {
			fmt.Printf("function %s (output parms %d) failed! \r\n", fun.name, fun.functype.NumOut())
			continue
		}

		fun.input[0] = fun.functype.In(0).String()
		fun.input[1] = fun.functype.In(1).String()
		fun.output = fun.functype.Out(0).String()

		if fun.output != "error" {
			fmt.Printf("function %s (output type %s) failed! \r\n", fun.name, fun.output)
			continue
		}

		s.symbol[fun.name] = fun

		fmt.Printf("Add Method: %s \r\n", fun.name)
	}
}

func (s *Server) Call(method string, req interface{}, rsp interface{}) error {

	s.lock.Lock()
	defer s.lock.Unlock()

	fun, b := s.symbol[method]
	if b == false {
		return errors.New("can not found " + method)
	}

	parms := make([]reflect.Value, 2)
	parms[0] = reflect.ValueOf(req)
	parms[1] = reflect.ValueOf(rsp)

	for i := 0; i < 2; i++ {
		if parms[i].Type().String() != fun.input[i] {
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

func (s *Server) Start() {

}
