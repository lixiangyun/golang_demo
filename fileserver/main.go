package main

import (
	"net/http"
)

func main() {
	h := http.FileServer(http.Dir("../"))
	http.ListenAndServeTLS(":8080", "server.crt", "server.key", h)
}
