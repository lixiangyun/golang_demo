package main

import (
	"fmt"
	"os"
)

func readtest() {
	userFile := "astaxie.txt"
	fl, err := os.Open(userFile)
	defer fl.Close()
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if 0 == n {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}

func writetest() {
	userFile := "astaxie.txt"
	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	for i := 0; i < 10; i++ {
		fout.WriteString("Just " + string('0'+byte(i)) + " test!\r\n")
		fout.Write([]byte("Just A Test!\r\n"))
	}
}

func deletefile() {
	userFile := "astaxie.txt"
	err := os.Remove(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
}

func main() {
	writetest()
	readtest()
	deletefile()
}
