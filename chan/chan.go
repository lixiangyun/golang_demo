package main

import (
	//"fmt"
	"time"
)
import "strconv"

type quelist struct {
	q    chan int
	id   int
	name string
}

var g_all *quelist

func func1(id int) {
	//one := g_all[id]

	//fmt.Println(g_all)

	//fmt.Println(one)
}

func func2() {
	var que quelist

	g_all = &que

	g_all.id = 0
	g_all.name = strconv.Itoa(0)
	g_all.q = make(chan int, 100)
}

func main() {

	for i := 0; i < 10000000; i++ {
		func2()
		func1(0)
		func2()
		func1(0)
	}

	time.Sleep(time.Duration(1) * time.Second)
}
