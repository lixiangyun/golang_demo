package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("c")
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		fmt.Println("d")
	}()
	f()
}

func send(conn chan struct{}) {

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("c")
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		fmt.Println("d")
	}()

	fmt.Println("c_1")
	conn <- struct{}{}
	fmt.Println("c_2")
	conn <- struct{}{}
	fmt.Println("c_3")

}

func f() {

	conn := make(chan struct{}, 1)

	go send(conn)

	time.Sleep(1 * time.Second)

	close(conn)

	fmt.Println("a")

	fmt.Println("b")
	fmt.Println("f")
}
