package main

import (
	"fmt"
	"io"
)
import "os"

type PathError struct {
	Opt  string
	Path string
	Err  error
}

func (e *PathError) Error() string {
	return e.Opt + " " + e.Path + ": " + e.Err.Error()
}

func Stat(name string) (fi PathError) {

	_, err := os.Stat(name)

	if nil != err {
		return PathError{"stat", name, err}
	} else {
		return PathError{"stat", name, nil}
	}
}

func CopyFile(dst, src string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if nil != err {
		return
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if nil != err {
		return
	}

	defer dstFile.Close()

	defer func() {
		fmt.Println("close done!")
	}()

	return io.Copy(dstFile, srcFile)
}

func main() {
	var j int = 5

	a := func() func() {
		var i int = 10

		return func() {
			fmt.Printf("i,j:%d,%d\n", i, j)
		}
	}()

	a()

	j *= 2

	a()

	err := Stat("./closure2.go")

	if nil != err.Err {
		fmt.Println(err.Opt, err.Path, err.Err.Error())
	} else {
		fmt.Println(err.Opt, err.Path, "ok!")
	}

	err = Stat("./closure.go")

	if nil != err.Err {
		fmt.Println(err.Opt, err.Path, err.Err.Error())
	} else {
		fmt.Println(err.Opt, err.Path, "ok!")
	}

	w, err2 := CopyFile("temp", "closure.go")

	if nil != err2 {
		fmt.Printf(err2.Error())
	} else {
		fmt.Println("w : ", w)
	}

}
