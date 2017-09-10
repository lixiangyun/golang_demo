package main

import (
	"encoding/json"
	"fmt"
)

type KvData struct {
	Flags bool   `json:"Flag"`
	Index int    `json:"LockIndex"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func main() {

	var kv KvData
	var kv2 KvData

	kv.Index = 1
	kv.Flags = true
	kv.Key = "abc"
	kv.Value = "hello world!"

	buf, err := json.Marshal(kv)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &kv2)
	if err != nil {
		return
	}

	fmt.Println(kv2)

	fmt.Println("json<body>:", string(buf))

	if kv2 == kv {
		fmt.Println("Json Marshal & Unmarshal success!")
	}

}
