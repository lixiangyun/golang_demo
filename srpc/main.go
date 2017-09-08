package main

import (
	"errors"
	"fmt"
	"golang_demo/srpc/srpc"
)

type SAVE struct {
}

func (s *SAVE) Add(a, b *int) error {

	if a == nil || b == nil {
		return errors.New("input parms error! ")
	}

	*b = (*a + 1)

	fmt.Println("call add ", *a, *b)

	return nil
}

func main() {

	var s SAVE

	server := srpc.NewServer("localhost:1234")

	server.Bind(&s)

	var a, b int
	a = 123
	b = 321

	err := server.Call("Add", &a, &b)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("a= ", a, " b=", b)
	}

	err = server.Call("Add", a, &b)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("a= ", a, " b=", b)
	}

	return
}
