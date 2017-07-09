package main

import "flag"
import "fmt"

type Strings interface {
	String() string
}

func PrtArgs(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Println("Int : ", arg)
		case string:
			fmt.Println("Str : ", arg)
		default:
			if v, ok := arg.(Strings); ok {
				val := v.String()
				fmt.Println("Strs : ", val)
			} else {
				fmt.Println("unkown type")
			}
		}
	}
}

func String() string {
	return "hello world"
}

func main() {

	flag.Parse()

	var v1 interface{} = 1
	var v2 interface{} = "abc"
	//var v3 interface{} = Strings{}

	var v3 Strings

	PrtArgs(flag.Args())

	PrtArgs(v1)
	PrtArgs(v2)
	PrtArgs(v3)

}
