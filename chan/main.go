package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var cnt int32

func proxy(queue chan int, wg *sync.WaitGroup)  {
	defer wg.Done()
	for  {
		_, close := <- queue
		if close == false {
			break
		}
		atomic.AddInt32(&cnt, 1)
	}
}

func main()  {
	wg := new(sync.WaitGroup)
	queue := make(chan int, 100)

	for i:=0; i<10; i++ {
		wg.Add(1)
		go proxy(queue, wg)
	}

	for i:=0; i<10000;i++ {
		queue <- i
	}

	close(queue)

	wg.Wait()

	fmt.Println(cnt)
}
