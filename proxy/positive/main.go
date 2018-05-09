package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	//"strings"
	"time"
)

var HTTP_PROXY string = ":808"

type proxy struct{}

func (*proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	fmt.Println(req)

	var transport *http.Transport

	transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second}

	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// step 1
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	// step 2
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// step 3
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func main() {

	http.ListenAndServe(HTTP_PROXY, &proxy{})
}
