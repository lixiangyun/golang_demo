package main

/*
HTTP/2 客户端例子
Author: XCL
Date: 2016-12-25
测试结果:
➜  client  : go run client.go
resp.Body:
 RequestURI: /v1
Protocol: HTTP/2.0
Hello V1
*/

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

func HttpClient1(url string) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := http.Client{Transport: tr}

	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancel()
	})

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("resp.Body:\n", string(body))
}

func HttpClient2(url string) {

	tr := &http2.Transport{
		AllowHTTP: true, //充许非加密的链接
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient := http.Client{Transport: tr}

	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancel()
	})

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println("resp.Body:\n", string(body))
}
