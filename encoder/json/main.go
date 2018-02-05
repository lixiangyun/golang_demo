package main

import (
	"encoding/json"
	"log"
)

type KvData struct {
	Flags bool   `json:"Flag"`
	Index int    `json:"LockIndex"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type ExtData struct {
	Key   string `json:"Key2"`
	Value string `json:"Value2"`
}

type CacheItem struct {
	Key    string `json:"key"`
	MaxAge int    `json:"cacheAge"`
	Value  int    `json:"cacheValue"`
}

func main() {

	var kv KvData
	var kv2 KvData

	kv.Index = 1
	kv.Flags = true
	kv.Key = "abc"
	kv.Value = "hello world!"

	var extkv ExtData
	var extkv2 ExtData

	extkv.Key = "name"
	extkv.Value = "Jack!"

	buf, err := json.Marshal(struct {
		*KvData
		*ExtData
	}{&kv, &extkv})
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(buf, &kv2)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(buf, &extkv2)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(kv2)

	log.Println("json<body>:", string(buf))

	if kv2 == kv {
		log.Println("Json Marshal & Unmarshal success 1!")
	}

	if extkv2 == extkv {
		log.Println("Json Marshal & Unmarshal success 2!")
	}

	item := &CacheItem{Key: "abcdef", MaxAge: 2, Value: 3}

	buf, err = json.Marshal(struct {
		*CacheItem

		// Omit bad keys
		OmitMaxAge int `json:"cacheAge,omitempty"`
		OmitValue  int `json:"cacheValue,omitempty"`

		// Add nice keys
		MaxAge int  `json:"max_age"`
		Value  *int `json:"value"`
	}{
		CacheItem: item,

		// Set the int by value:
		MaxAge: item.MaxAge,

		// Set the nested struct by reference, avoid making a copy:
		Value: &item.Value,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("json<body>:", string(buf))

}
