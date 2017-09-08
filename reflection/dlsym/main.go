package main

import (
	"fmt"
	"reflect"
)

//定义控制器函数Map类型，便于后续快捷使用
type ControllerMapsType map[string]reflect.Value

//声明控制器函数Map类型变量
var ControllerMaps ControllerMapsType

//定义路由器结构类型
type Routers struct {
}

//为路由器结构附加功能控制器函数，值传递
func (this *Routers) Login(msg string) {
	fmt.Println("Login:", msg)
}

//为路由器结构附加功能控制器函数，引用传递
func (this *Routers) ChangeName(msg *string) {
	fmt.Println("ChangeName:", *msg)
	*msg = *msg + " Changed"
}

func main() {
	var ruTest Routers

	crMap := make(ControllerMapsType, 0)
	//创建反射变量，注意这里需要传入ruTest变量的地址；
	//不传入地址就只能反射Routers静态定义的方法
	vf := reflect.ValueOf(&ruTest)
	vft := vf.Type()

	//读取方法数量
	mNum := vf.NumMethod()
	fmt.Println("NumMethod:", mNum)
	//遍历路由器的方法，并将其存入控制器映射变量中
	for i := 0; i < mNum; i++ {
		mName := vft.Method(i).Name
		fmt.Println("index:", i, " MethodName:", mName)
		crMap[mName] = vf.Method(i) //<<<
	}

	//演示
	testStr := "Hello Go"
	//创建带调用方法时需要传入的参数列表
	parms := []reflect.Value{reflect.ValueOf(testStr)}
	//使用方法名字符串调用指定方法
	crMap["Login"].Call(parms)

	//创建带调用方法时需要传入的参数列表
	parms = []reflect.Value{reflect.ValueOf(&testStr)}
	//使用方法名字符串调用指定方法
	crMap["ChangeName"].Call(parms)
	//可见，testStr的值已经被修改了
	fmt.Println("testStr:", testStr)
}
