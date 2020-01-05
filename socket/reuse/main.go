package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gogf/greuse"
)

// We can create two processes with this code.
// Do some requests, then watch the output of the console.
func main() {
	listener, err := greuse.Listen("tcp", ":8881")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	server := &http.Server{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "gid: %d, pid: %d\n", os.Getgid(), os.Getpid())
	})

	panic(server.Serve(listener))
}
