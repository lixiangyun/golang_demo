package main

import (
	"errors"
	"fmt"
	"golang_demo/srpc/srpc"
)

type SAVE struct {
}

func (s *SAVE) Add(a, b *int) error {

	if *a > 100 {
		return errors.New("input error!")
	}

	*b = (*a + 1)

	fmt.Println("call add ", *a, *b)

	return nil
}

func main() {

	var s SAVE

	server := srpc.NewServer("localhost:1234")

	server.BindMethod(&s)

	var a, b int
	a = 1
	b = 2
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

	a = 123
	b = 2
	err = server.Call("Add", &a, &b)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("a= ", a, " b=", b)
	}
	return
}
