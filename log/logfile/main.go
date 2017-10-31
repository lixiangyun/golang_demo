package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	file, err := os.OpenFile("runlog.txt", os.O_WRONLY, 0)
	if err != nil {
		file, err = os.Create("runlog.txt")
		if err != nil {
			fmt.Println("create file error!")
			return
		}
		fmt.Println("create log file ")
	} else {
		fileinfo, err := file.Stat()
		if err == nil {
			file.Seek(fileinfo.Size(), 0)
			fmt.Println("append log to file ", file.Name())
		}
	}

	defer file.Close()

	log.SetOutput(file)

	log.Println("helloworld!")

	return
}
