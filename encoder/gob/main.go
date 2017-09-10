package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

type User struct {
	Id   int
	Name string
}

func (this *User) Say() string {
	return this.Name + ` hello world ! `
}

func encoder() {
	file, err := os.Create("mygo.gob")
	if err != nil {
		fmt.Println(err)
		return
	}
	user := User{Id: 1, Name: "Mike"}
	user2 := User{Id: 3, Name: "Jack"}
	u := []User{user, user2}

	enc := gob.NewEncoder(file)

	err2 := enc.Encode(u)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
}

func decoder() {

	var u []User
	file, err := os.Open("mygo.gob")
	if err != nil {
		fmt.Println(err)
	}
	dec := gob.NewDecoder(file)
	err2 := dec.Decode(&u)

	if err2 != nil {
		fmt.Println(err2)
		return
	}

	for _, user := range u {
		fmt.Println(user.Id)
		fmt.Println(user.Say())
	}
}

func main() {
	encoder()
	decoder()
}
