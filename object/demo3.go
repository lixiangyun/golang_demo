package main

import "fmt"

type Integer int

func (a Integer) Less(b Integer) bool {
	return a < b
}

func main() {
	var v1 Integer = 1

	if v1.Less(2) {
		fmt.Println(v1, "Less 2")
	}

}
