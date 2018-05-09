package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []Server `xml:"server"`
}

func main() {
	file, err := os.Open("test.xml") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	var v Servers

	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Println(v)

	data2, err := xml.Marshal(v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	err = ioutil.WriteFile("test2.xml", data2, os.ModePerm)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	fmt.Println(string(data2))
}
