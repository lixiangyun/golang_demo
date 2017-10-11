package main

import (
	"fmt"
)

// 打印字符
func main() {
	var abc [12]byte

	cdf := abc[4:]

	for i, _ := range abc {
		abc[i] = byte(i)
	}

	fmt.Println(abc, &cdf)
}
