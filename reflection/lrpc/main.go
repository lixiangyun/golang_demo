package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"runtime/debug"
)

// 函数反射（包括参数类型）
type funcinfo struct {
	funvalue reflect.Value
	input    []reflect.Type
	output   []reflect.Type
}

// 符号表
var funtable map[string]*funcinfo

// 报文序列化
func CodePacket(req interface{}) ([]byte, error) {
	iobuf := new(bytes.Buffer)

	enc := gob.NewEncoder(iobuf)

	err := enc.Encode(req)
	if err != nil {
		debug.PrintStack()
		return nil, err
	}

	return iobuf.Bytes(), nil
}

// 报文反序列化
func DecodePacket(buf []byte, rsp interface{}) error {
	iobuf := bytes.NewReader(buf)
	denc := gob.NewDecoder(iobuf)
	err := denc.Decode(rsp)

	if err != nil {
		debug.PrintStack()
	}

	return err
}

func AddMethod(pthis interface{}) {

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

		fun.funvalue = vfun.Method(i)
		functype := vfun.Method(i).Type()

		inputnum := functype.NumIn()
		fmt.Println("inputnum: ", inputnum)

		fun.input = make([]reflect.Type, inputnum)
		for i := 0; i < inputnum; i++ {
			fun.input[i] = functype.In(i)
			fmt.Println("inputtype: ", fun.input[i].String())
		}

		outputnum := functype.NumOut()
		fmt.Println("outputnum: ", outputnum)

		fun.output = make([]reflect.Type, outputnum)
		for i := 0; i < outputnum; i++ {
			fun.output[i] = functype.Out(i)
			fmt.Println("outputtype: ", fun.input[i].String())
		}

		funname := vtype.Method(i).Name

		funtable[funname] = &fun

		fmt.Printf("Add Method: %s \r\n", funname)
	}
}

type SAVE struct {
	Abc uint32
}

// 构造验证的方法
func (s *SAVE) Add(a uint32, b uint32) uint32 {

	c := (a + b)

	fmt.Println("call Add ", a, b, c, s.Abc)

	return c
}

func (s *SAVE) Sub(a uint32, b uint32) uint32 {

	c := (a - b)

	fmt.Println("call sub ", a, b, c, s.Abc)

	return c
}

// 本地符号调用方法
func CallMethod(funname string, input []interface{}, output []interface{}) {

	funinfo, b := funtable[funname]
	if b == false {
		fmt.Println("func not found! ", funname)
		return
	}

	if len(funinfo.input) != len(input) {
		fmt.Println(len(funinfo.input), len(input))
		fmt.Println("input parm num error!")
		return
	}

	if len(funinfo.output) != len(output) {
		fmt.Println(len(funinfo.output), len(output))
		fmt.Println("output parm num error!")
		return
	}

	inputnum := len(input)
	inputtable := make([]reflect.Value, inputnum)
	for i := 0; i < inputnum; i++ {
		inputtable[i] = reflect.ValueOf(input[i])

		if inputtable[i].Type().String() != funinfo.input[i].String() {

			fmt.Println(inputtable[i].Type().String(), funinfo.input[i].String())
			debug.PrintStack()
			return
		}
	}

	outputnum := len(output)
	outputtable := make([]reflect.Value, outputnum)
	for i := 0; i < outputnum; i++ {
		outputtable[i] = reflect.ValueOf(output[i])
	}

	returntable := funinfo.funvalue.Call(inputtable)

	returnnum := len(returntable)
	for i := 0; i < returnnum; i++ {

		fmt.Println("return ", returntable[i].Interface())

		temp := reflect.Indirect(outputtable[i])
		temp.Set(returntable[i])
	}

	return
}

func main() {

	funtable = make(map[string]*funcinfo)

	var temp SAVE
	AddMethod(&temp)

	// 测试用例
	input := make([]interface{}, 2)
	output := make([]interface{}, 1)

	var a, b, c uint32

	a = 3
	b = 2

	temp.Abc = 1000

	// 测试Add方法
	input[0] = a
	input[1] = b
	output[0] = &c
	CallMethod("Add", input, output)
	fmt.Println(a, b, c)

	a = 10
	b = 4

	temp.Abc = 2000

	// 测试sub方法
	input[0] = a
	input[1] = b
	output[0] = &c
	CallMethod("Sub", input, output)
	fmt.Println(a, b, c)
}
