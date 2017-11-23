package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func strcat(s1 string, cut string) string {
	begin := strings.Index(s1, cut)
	if begin != -1 {
		return s1[0:begin] + s1[begin+len(cut):]
	} else {
		return s1
	}
}

func strcatall(s1 string, cut string) string {
	newbody := s1
	for {
		length := len(newbody)
		newbody = strcat(newbody, cut)
		if length == len(newbody) {
			break
		}
	}
	return newbody
}

func parse(body string) string {

	newbody := strcatall(body, "</p>")
	newbody = strcatall(newbody, "<p>")
	newbody = strcatall(newbody, "\t")
	newbody = strcatall(newbody, "\r")
	newbody = strcatall(newbody, "\n")
	newbody = strcatall(newbody, "&rsquo")
	newbody = strcatall(newbody, "&lsquo")

	return newbody
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

func getcvebody(buf []byte) string {
	body := string(buf[:])

	head := "<p data-testid=\"vuln-description\">"
	tail := "</p>"

	begin := strings.Index(body, head)
	if begin != -1 {
		end := strings.Index(body[begin:], tail)
		if end != -1 {
			cutbody := body[begin+len(head) : begin+end]
			return parse(cutbody)
		}
	}

	return ""
}

func getcvelist(filename string) ([]string, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var body string

	for {
		var buf [128]byte
		cnt, err := file.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		body += string(buf[:cnt])
	}

	return strings.Split(body, "\r\n"), nil
}

func getcveall(in chan string, out chan string) {

	for {
		oneline, b := <-in
		if b == false {
			log.Println("close go routine.")
			return
		}

		log.Println("get cve body : ", oneline)

		vlist := strings.Split(oneline, ",")
		var output string

		for _, cve := range vlist {
			if strings.Index(cve, "CVE") == -1 {
				continue
			}

			buf, err := gethtmlbody("https://nvd.nist.gov/vuln/detail/"+cve, "http://10.177.16.233:808")
			if err != nil {
				log.Println("gethtmlbody : ", cve, err.Error())
				output += fmt.Sprintf("%s Get Body Failed! \r\n", cve)
				continue
			}
			body := getcvebody(buf)
			output += fmt.Sprintf("%s\t%s \r\n", cve, body)
		}

		out <- output
	}
}

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Println("input name failed!")
		return
	}

	list, err := getcvelist(args[1])
	if err != nil {
		log.Println("getcvelist : ", err.Error())
		return
	}

	log.Println("get total : ", len(list))

	in := make(chan string, 10)
	out := make(chan string, 10)

	for i := 0; i < 5; i++ {
		go getcveall(in, out)
	}

	go func() {
		for _, v := range list {
			in <- v
		}
	}()

	var num int

	for {
		body := <-out
		fmt.Println(body)
		num++
		if num >= len(list) {
			break
		}
	}

	close(in)
}
