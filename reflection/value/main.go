package main

import (
	"fmt"
	"reflect"
)

func hello(a *int, b int) {

	value1 := reflect.ValueOf(a)
	//value2 := reflect.ValueOf(b)

	vtype1 := reflect.TypeOf(a)
	vtype2 := reflect.TypeOf(b)

	fmt.Println(vtype1)
	fmt.Println(vtype2)

	value11 := reflect.Indirect(value1)
	if value11 != value1 {
		fmt.Println(value11.Type())
		fmt.Println(value1.Type())
	}

	aa := reflect.New(vtype1)
	bb := reflect.New(vtype2)

	fmt.Println(aa.Type())
	fmt.Println(bb.Type())

	a1 := reflect.Indirect(aa)
	b1 := reflect.Indirect(bb)

	b1.SetInt(2)

	fmt.Println(a1.Type())
	fmt.Println(b1.Type())

	fmt.Println(a1.Interface())
	fmt.Println(b1.Interface())

}

func main() {

	var a, b int

	a = 1
	b = 2

	hello(&a, b)
}
