package main

import "fmt"

// 动态参数传递实例
func MyPrintf(args ...interface{}) {
	for _, v := range args {
		switch v.(type) {
		case int:
			fmt.Println(v, "is an int value.")
		case string:
			fmt.Println(v, "is an string value.")
		case int64:
			fmt.Println(v, "is an int64 value.")
		default:
			fmt.Println(v, "is an unkown type.")
		}
	}
}

func main() {
	var v1 int = 1
	var v2 string = "helloworld"
	var v3 int64 = 0xfffffffffffffff
	var v4 float32 = 0.1

	MyPrintf(v1, v2, v3, v4)

	f := func(x, y int) int {
		return x + y
	}

	v := f(1, 2)

	fmt.Println("v : ", v)

	func(ch int) {
		fmt.Println("v2: ", ch)
	}(v)

}
