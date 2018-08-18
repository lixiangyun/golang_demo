package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func ClientTlsConfig() *tls.Config {

	//这里读取的是根证书
	buf, err := ioutil.ReadFile(CA_FILE)
	if err != nil {
		log.Fatalln(err.Error())
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(buf)

	//加载客户端证书
	//这里加载的是服务端签发的
	cert, err := tls.LoadX509KeyPair(CLIENT_CERT, CLIENT_KEY)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &tls.Config{
		ServerName:   "127.0.0.1",
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
	}
}

func newTransport() http.RoundTripper {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       ClientTlsConfig(),
	}
}

func newhttpclient() *http.Client {
	return &http.Client{
		Transport: newTransport(),
		Timeout:   10 * time.Second,
	}
}

func HttpClient(addr string) {

	httpclient := newhttpclient()

	request, err := http.NewRequest("GET",
		"https://"+addr+"/home",
		strings.NewReader("hello world! form https client!"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("request : ", request)

	resp, err := httpclient.Do(request)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("response : ", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("body : ", string(body))
}
