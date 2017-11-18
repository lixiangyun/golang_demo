package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func hashfile() {
	TestFile := "123.txt"
	file, err := os.Open(TestFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	md5Inst := md5.New()
	io.Copy(md5Inst, file)
	Result := md5Inst.Sum([]byte(""))
	fmt.Printf("%x\r\n", Result)

	ShalInst := sha1.New()
	io.Copy(ShalInst, file)
	Result = ShalInst.Sum([]byte(""))
	fmt.Printf("%x\r\n", Result)

}

func main() {
	TestString := "Hi,pandaman!"

	Md5Inst := md5.New()
	Md5Inst.Write([]byte(TestString))
	Result := Md5Inst.Sum([]byte(""))
	fmt.Printf("%x\r\n", Result)

	ShalInst := sha1.New()
	ShalInst.Write([]byte(TestString))
	Result = ShalInst.Sum([]byte(""))
	fmt.Printf("%x\r\n", Result)

	hashfile()
}
