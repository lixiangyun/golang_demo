package main

import "fmt"

type Base struct {
	Name string
	id   int
}

func (b *Base) GetName() string {
	return b.Name
}

func (b *Base) GetId() int {
	return b.id
}

type Jack struct {
	Base
	add string
}

func (j *Jack) GetName() string {
	return j.Name + " " + j.add
}

func main() {

	b1 := Base{"bar", 123}
	j1 := Jack{Base{"Jak", 333}, "home 001"}

	fmt.Println(b1, b1.GetName())
	fmt.Println(j1, j1.GetName())

}
