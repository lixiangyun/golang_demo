package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	HELP bool
	DIR  string
	CONTENT string
	Separator string
)

func  init()  {
	Separator = fmt.Sprintf("%c",os.PathSeparator)
	flag.BoolVar(&HELP,"help",false,"This help.")
	flag.StringVar(&DIR,"dir","."+Separator,"The directory path to be searched.")
	flag.StringVar(&CONTENT,"ctx","","Search for matched content,")
}

func matchFile(filename string,ctx string)  {

	body , err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s open failed! (%s)\n",filename,err.Error())
		return
	}

	buffer :=	bytes.NewBuffer(body)

	line := 1
	for {
		oneline, err:= buffer.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("%s open failed! (%s)\n",filename,err.Error())
			}
			return
		}

		if -1 != strings.Index(oneline,ctx) {

			fmt.Printf("[%s:%d] %s \r\n",filename,line,oneline)
		}

		line++
	}
}

func listFile(dir string) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			listFile(dir + Separator + file.Name())
		} else {
			filename := file.Name()
			matchFile(dir + Separator + filename,CONTENT)
		}
	}
}

func main() {

	flag.Parse()
	if HELP || CONTENT == ""{
		flag.Usage()
		os.Exit(-1)
	}

	listFile(DIR)
}
