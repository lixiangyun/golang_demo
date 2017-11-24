package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	help     bool
	crtcheck bool
	filename string
	proxy    string
	path     string
	timeout  int
)

func usage() {
	fmt.Fprintf(os.Stderr,
		`wget version: 1.0
Usage: wget [-h] [-c] [-t second] [-o filename] [-p proxy] [-u url]

Options:
`)
	flag.PrintDefaults()
}

func init() {

	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&crtcheck, "c", false, "validate the server's certificate")
	flag.StringVar(&filename, "o", "", "output the body filename")
	flag.StringVar(&proxy, "p", "", "the proxy server url")
	flag.StringVar(&path, "u", "", "the URL need to get")
	flag.IntVar(&timeout, "t", 0, "timeout for get by tcp connect")

	flag.Usage = usage

}

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

func gethtmlbody() ([]byte, error) {

	var transport *http.Transport

	transport = new(http.Transport)

	if proxy != "" {
		proxyfunc := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxy)
		}
		transport.Proxy = proxyfunc
	}

	if crtcheck == false {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: transport, Timeout: time.Duration(timeout) * time.Second}

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

func SaveFile(body []byte) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Errorf("create file failed!", filename)
		return
	}

	defer file.Close()

	file.Write(body)
}

func main() {

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	body, err := gethtmlbody()
	if err != nil {
		log.Println(err.Error())
		return
	}

	if filename != "" {

	} else {
		fmt.Fprintln(os.Stdout, string(body[:]))
	}
}
