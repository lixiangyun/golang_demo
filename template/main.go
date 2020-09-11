package main

import (
	"os"
	"log"
	"text/template"
)

type ProxyCfg struct {
	Https    int
	Name     string
	CertFile string
	CertKey  string
	Backend  string
	LogDir   string
}

type Config struct {
	Redirect bool
	LogDir   string
	Proxy  []ProxyCfg
}

func main() {
	http1 := ProxyCfg{
		Https: 881,
		Name: "test1.abc",
		CertKey: "/etc/key1.file",
		CertFile: "/etc/cert1.file",
		Backend: "http://192.168.1.1:18181",
	}

	http2 := ProxyCfg{
		Https: 443,
		Name: "test2.abc",
		CertKey: "/etc/key2.file",
		CertFile: "/etc/cert2.file",
		Backend: "http://192.168.1.2:18283",
	}

	ctx := &Config{
		Redirect: true,
		LogDir: "/home/https",
		Proxy: []ProxyCfg{
			http1, http2,
		},
	}

	t, err := template.ParseFiles("./nginx.conf.template")
	if err != nil {
		log.Println(err.Error())
		return
	}

	file, err := os.Create("nginx.conf")
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer file.Close()

	err = t.Execute(file, ctx)
	if err != nil {
		log.Println("Executing template:", err)
	}
}