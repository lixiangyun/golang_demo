package main

import (
	"fmt"
	"runtime"
	"sync"
)

const (
	MAX_CNT = 10000
)

var counter int = 0

func Count(lock *sync.Mutex) {
	for i := 0; i < MAX_CNT; i++ {
		lock.Lock()
		counter++
		lock.Unlock()
	}
}

func main() {
	lock := sync.Mutex{}

	for i := 0; i < 10; i++ {
		go Count(&lock)
	}

	for {
		lock.Lock()

		c := counter

		lock.Unlock()

		runtime.Gosched()

		if c >= MAX_CNT*10 {
			break
		}
	}

	fmt.Println("cnt : ", counter)
}
