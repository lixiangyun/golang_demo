package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	myfolder := os.Args[1]
	listFile(myfolder)
}

func listFile(myfolder string) {
	files, _ := ioutil.ReadDir(myfolder)
	for _, file := range files {
		if file.IsDir() {
			listFile(myfolder + "/" + file.Name())
		} else {
			filename := file.Name()

			if strings.Index(filename, ".md") != -1 {
				fmt.Println(myfolder + "/" + filename)
			}
		}
	}
}
