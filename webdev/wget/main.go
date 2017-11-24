package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func gethtmlbody(path string, proxy string) ([]byte, error) {

	var transport *http.Transport

	if proxy != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}

		if strings.Index(path, "https") != -1 {
			transport = &http.Transport{Proxy: proxy,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		} else {
			transport = &http.Transport{Proxy: proxy}
		}
	}

	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	var err error
	var resp *http.Response

	for i := 0; i < 100; i++ {
		resp, err = client.Get(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	buf, err := readFully(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func main() {

	body, err := gethtmlbody("https://www.google.com", "http://127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(string(body[:]))
}
