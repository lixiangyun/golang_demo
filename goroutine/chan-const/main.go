package main

import (
	"fmt"
	"time"
)

func RecvMsg(ch <-chan int) {
	for v := range ch {
		fmt.Println("v : ", v)
	}
}

func SendMsg(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
}

func main() {
	ch1 := make(chan int)

	var ch2 <-chan int
	var ch3 chan<- int

	ch2 = ch1
	ch3 = ch1

	go RecvMsg(ch2)
	go SendMsg(ch3)

	time.Sleep(time.Duration(1) * time.Second)

}
