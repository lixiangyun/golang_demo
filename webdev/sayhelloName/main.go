package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		sayHelloName(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form)
	fmt.Println("Path:", r.URL.Path)
	fmt.Println("Scheme:", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("Key:", k)
		fmt.Println("Val:", strings.Join(v, ""))
	}

	//fmt.Println(r)

	fmt.Fprintf(w, "Hello astaxie!")
}

func main() {
	mux := &MyMux{}
	http.HandleFunc("/", sayHelloName)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
