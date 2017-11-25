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

/* 例子：
 * wget.exe -c -o home.html -p http://127.0.0.1:8080 -t 10 -u http://www.google.com
 */

var (
	help     bool   // help标志，查看使用帮助
	crtcheck bool   // 是否检查证书，默认不检查
	filename string // 输出的文件名称
	proxy    string // http、https 访问代理
	path     string // 需要get的url路径，需要包含完整路径：http、https://xxx
	timeout  int    // 超时时长，单位second，默认10秒
	retry    int    // 重试次数，默认10次
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

	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&crtcheck, "c", false, "validate the server's certificate")
	flag.StringVar(&filename, "o", "home.html", "output the body filename")
	flag.StringVar(&proxy, "p", "", "the proxy server url")
	flag.StringVar(&path, "u", "", "the URL need to get")
	flag.IntVar(&timeout, "t", 10, "timeout for get by tcp connect")
	flag.IntVar(&retry, "r", 10, "retry times")

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
		SaveFile(body)
	} else {
		fmt.Fprintln(os.Stdout, string(body[:]))
	}
}
