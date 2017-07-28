package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func readFully(conn io.ReadCloser) ([]byte, error) {

	result := bytes.NewBuffer(nil)

	var buf [512]byte

	for {
		n, err := conn.Read(buf[0:])

		result.Write(buf[0:n])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}

	return result.Bytes(), nil
}

func main() {
	resp, err := http.Get("http://www.baidu.com/")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	buf, err := readFully(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Len:", len(buf))
	fmt.Println("Body:", string(buf))
}
