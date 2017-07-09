package main

import "fmt"
import "one"
import "two"

type Integer int

func (a Integer) Min(b Integer) bool {
	return a < b
}

func (a *Integer) Add(b Integer) Integer {
	*a += b
	return *a
}

func (a *Integer) Sub(b Integer) Integer {
	*a -= b
	return *a
}

type IntMath interface {
	Min(b Integer) bool
	Add(b Integer) Integer
	Sub(b Integer) Integer
}

type Lesser interface {
	Min(b Integer) bool
}

func main() {
	var a Integer = 1
	var b Integer = 2

	fmt.Println("min : ", a.Min(b))
	fmt.Println("Add : ", a.Add(b))
	fmt.Println("Sub : ", b.Sub(a))

	fmt.Println(a, b)

	var c IntMath = &a

	c.Sub(b)

	fmt.Println(a, b, c)

	var c2 Lesser = &a
	var c3 Lesser = a

	c2.Min(b)
	c3.Min(b)

	fmt.Println(a, b, c, c2, c3.Min(b))

}
