package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

var (
	LISTEN_ADDR   string
	REDIRECt_ADDR string

	DEBUG bool
	help  bool
)

func init() {
	flag.StringVar(&LISTEN_ADDR, "in", "", "listen addr for http/https proxy.")
	flag.StringVar(&REDIRECt_ADDR, "out", "", "redirect to addr for http/https proxy.")
	flag.BoolVar(&DEBUG, "debug", false, "debug mode.")
	flag.BoolVar(&help, "h", false, "this help.")
}

type HttpProxy struct {
	Addr      string
	RedirAddr string
	Svc       *http.Server
	GoCnt     int
	Que       chan *HttpRequest
	Wait      sync.WaitGroup
	Stop      chan struct{}
}

type HttpRsponse struct {
	status int
	header http.Header
	body   []byte
	err    error
}

var requestNum int32

type HttpRequest struct {
	num    int32
	addr   string
	url    string
	method string
	header http.Header
	body   []byte
	rsp    chan *HttpRsponse
}

func newHttpClient() *http.Client {

	tr := &http2.Transport{
		AllowHTTP: true, //充许非加密的链接
		DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(netw, addr)
		},
	}

	return &http.Client{Transport: tr,
		Timeout: 10 * time.Second}
}

func (h *HttpProxy) Process() {

	defer h.Wait.Done()
	httpclient := newHttpClient()

	for {
		select {
		case proxyreq := <-h.Que:
			{
				proxyrsp := new(HttpRsponse)

				request, err := http.NewRequest(proxyreq.method,
					proxyreq.url,
					bytes.NewBuffer(proxyreq.body))
				if err != nil {
					proxyrsp.err = err
					proxyrsp.status = http.StatusInternalServerError

					proxyreq.rsp <- proxyrsp
					continue
				}

				for key, value := range proxyreq.header {
					for _, v := range value {
						request.Header.Add(key, v)
					}
				}

				resp, err := httpclient.Do(request)
				if err != nil {
					proxyrsp.err = err
					proxyrsp.status = http.StatusInternalServerError

					proxyreq.rsp <- proxyrsp
					continue
				} else {
					proxyrsp.status = resp.StatusCode
					proxyrsp.header = resp.Header
				}

				proxyrsp.body, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					proxyrsp.err = err
					proxyrsp.status = http.StatusInternalServerError

					proxyreq.rsp <- proxyrsp
					continue
				}
				resp.Body.Close()

				proxyreq.rsp <- proxyrsp
			}
		case <-h.Stop:
			{
				return
			}
		}
	}
}

func (h *HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	var err error

	defer req.Body.Close()

	redirect := h.RedirAddr

	// step 1
	proxyreq := new(HttpRequest)
	proxyreq.num = atomic.AddInt32(&requestNum, 1)
	proxyreq.addr = redirect
	proxyreq.method = req.Method
	proxyreq.header = req.Header
	proxyreq.rsp = make(chan *HttpRsponse, 1)

	proxyreq.url = "http://" + redirect + req.URL.RequestURI()

	proxyreq.body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if DEBUG {
		headers := fmt.Sprintf("\r\nHeader:\r\n")
		for key, value := range proxyreq.header {
			headers += fmt.Sprintf("\t%s:%v\r\n", key, value)
		}
		var body string
		if len(proxyreq.body) > 0 {
			body = fmt.Sprintf("Body:%s\r\n", string(proxyreq.body))
		}
		log.Printf("[%d]Request Method:%s\r\nURL:%s%s%s\r\n",
			proxyreq.num, proxyreq.method, proxyreq.url, headers, body)
	}

	h.Que <- proxyreq
	proxyrsp := <-proxyreq.rsp

	if DEBUG {
		headers := fmt.Sprintf("\r\nHeader:\r\n")
		for key, value := range proxyrsp.header {
			headers += fmt.Sprintf("\t%s:%v\r\n", key, value)
		}
		var body string
		if len(proxyrsp.body) > 0 {
			body = fmt.Sprintf("Body:%s\r\n", string(proxyrsp.body))
		}
		log.Printf("[%d]Response Code:%d%s%s\r\n",
			proxyreq.num, proxyrsp.status, headers, body)
	}

	// step 2
	if proxyrsp.err != nil {
		log.Println(proxyrsp.err.Error())
		http.Error(rw, proxyrsp.err.Error(), http.StatusInternalServerError)
		return
	}

	// step 3
	for key, value := range proxyrsp.header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(proxyrsp.status)
	rw.Write(proxyrsp.body)
}

func NewHttp2Proxy(addr string, redirAddr string) error {
	proxy := new(HttpProxy)

	proxy.Addr = addr
	proxy.RedirAddr = redirAddr

	list, err := net.Listen("tcp", proxy.Addr)
	if err != nil {
		log.Println("http listen failed!", err.Error())
		return err
	}
	defer list.Close()

	log.Printf("Http Proxy Listen %s\r\n", addr)

	Svc2 := &http2.Server{}

	proxy.Svc = &http.Server{
		Handler:      proxy,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second}

	proxy.Que = make(chan *HttpRequest, 1000)
	proxy.Stop = make(chan struct{}, proxy.GoCnt)

	for i := 0; i < 100; i++ {
		go proxy.Process()
	}

	if DEBUG {
		http2.VerboseLogs = true
	}

	http2.ConfigureServer(proxy.Svc, Svc2)

	for {
		conn, err := list.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}
		go Svc2.ServeConn(conn, &http2.ServeConnOpts{BaseConfig: proxy.Svc})
	}
}

func main() {

	flag.Parse()

	if help || LISTEN_ADDR == "" || REDIRECt_ADDR == "" {
		flag.Usage()
		return
	}

	log.Printf("Listen   At [%s]\r\n", LISTEN_ADDR)
	log.Printf("Redirect To [%s]\r\n", REDIRECt_ADDR)

	NewHttp2Proxy(LISTEN_ADDR, REDIRECt_ADDR)
}
