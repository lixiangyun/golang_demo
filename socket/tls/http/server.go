package main

import (
	"net/http"

	"crypto/tls"
	"crypto/x509"

	"io/ioutil"
	"log"
	"net"
)

type HttpProxy struct {
	Svc *http.Server
}

func (h *HttpProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	var err error

	defer req.Body.Close()

	log.Println("request : ", req.URL.RequestURI(), req.Method)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(body)
}

func ServerTlsConfig() *tls.Config {

	//这里读取的是根证书
	buf, err := ioutil.ReadFile(CA_FILE)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(buf)

	//加载服务端证书
	cert, err := tls.LoadX509KeyPair(SERVER_CERT, SERVER_KEY)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}
}

func HttpServer(addr string) {

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("https listen failed!", err.Error())
		return
	}
	log.Printf("https Proxy Listen: %s!\r\n", addr)

	httpServer := &http.Server{Handler: &HttpProxy{}, TLSConfig: ServerTlsConfig()}

	httpServer.ServeTLS(lis, "", "")
}
