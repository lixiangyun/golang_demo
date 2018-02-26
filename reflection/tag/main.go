package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name   string "user name"
	Passwd string "user password"
}

type Ser struct {
	f string `json:"Flag" xml:"flag"`
}

func main() {

	u := User{Name: "Jack", Passwd: "123456"}
	ref := reflect.TypeOf(&u)
	s := ref.Elem()

	for i := 0; i < s.NumField(); i++ {
		fmt.Println(s.Field(i).Tag)
	}

	flag := Ser{}
	ref1 := reflect.TypeOf(&flag)
	ss := ref1.Elem()

	for i := 0; i < ss.NumField(); i++ {
		fmt.Println(ss.Field(i).Tag.Get("json"),
			ss.Field(i).Tag.Get("xml"),
			ss.Field(i).Tag.Get("123"),
			ss.Field(i).Tag)
	}
}
