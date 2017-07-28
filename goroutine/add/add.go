package main

import "fmt"

func Add(x, y int) {
	z := x + y
	fmt.Println("z : ", z)
}

func main() {
	for i := 0; i < 10; i++ {
		go Add(i, i)
	}
}
