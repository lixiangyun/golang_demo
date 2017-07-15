package main

import "fmt"

func Count(id int, que chan int) {
	fmt.Println("Counting")
	que <- id
}

func main() {
	var que = make([]chan int, 10)

	for i := 0; i < 10; i++ {
		que[i] = make(chan int)
		go Count(i, que[i])
	}

	for _, ch := range que {
		id, _ := <-ch
		fmt.Println("RecvCnt : ", id)
	}
	
	ch1 := make(chan int ,1)
	
	for {
		switch {
			case ch1<-1:
			fmt.Println()
			case ch1<-0:
			fmt.Println()
			default:
			fmt.Println()
		}
		
		i:= <- ch1
		
		fmt.Println(" i : ", i)
	}
}
