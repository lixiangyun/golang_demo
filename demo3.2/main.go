package main

import "fmt"

type Rect struct {
	x, y float32
	w, h float32
}

func (r *Rect) Area() float32 {
	return r.w * r.h
}

func NewRect(x, y, w, h float32) *Rect {
	return &Rect{x, y, w, h}
}

func main() {
	r0 := new(Rect)
	r1 := &Rect{}
	r2 := &Rect{1, 2, 3, 4}
	r3 := &Rect{w: 5, h: 6}
	r4 := NewRect(6, 6, 6, 6)

	fmt.Println("r0 : ", r0, r0.Area())
	fmt.Println("r1 : ", r1, r1.Area())
	fmt.Println("r2 : ", r2, r2.Area())
	fmt.Println("r3 : ", r3, r3.Area())
	fmt.Println("r4 : ", r4, r4.Area())
}
