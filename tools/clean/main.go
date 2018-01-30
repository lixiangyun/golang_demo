package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func OpenFile(filename string) ([]byte, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	var body []byte
	for {
		var tmp [128]byte
		cnt, err := fd.Read(tmp[:])
		if err != nil {
			if io.EOF == err {
				break
			}
			return nil, err
		}
		body = append(body, tmp[:cnt]...)
	}

	return body, nil
}

func SaveFile(filename string, body []byte) error {

	var fd *os.File
	var err error

	for {
		fd, err = os.Create(filename)
		if err == nil {
			break
		}
		if os.ErrExist != err {
			return err
		}
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}

	defer fd.Close()

	cnt := 0
	for {
		num, err := fd.Write(body[cnt:])
		if err != nil {
			return err
		}
		cnt += num
		if cnt == len(body) {
			break
		}
	}

	return nil
}

func swap(input, oldvalue, newvalue string) string {
	for {
		idx := strings.Index(input, oldvalue)
		if idx == -1 {
			return input
		}
		input = fmt.Sprintf("%s%s%s", input[:idx], newvalue, input[idx+len(oldvalue):])
	}
	return input
}

func addchar(input string) string {
	output := make([]rune, 0)
	var chars bool
	var begin bool
	for _, v := range input {
		if chars == false {
			if unicode.IsLower(v) {
				output = append(output, '`')
				chars = true
				begin = true
			} else if unicode.IsUpper(v) {
				chars = true
			}
		} else {
			if !unicode.IsLower(v) && !unicode.IsUpper(v) {
				if v != '_' && v != '-' {
					chars = false
					if begin {
						output = append(output, '`')
						begin = false
					}
				}
			}
		}
		output = append(output, v)
	}
	return string(output)
}

func process(input, output string) {
	body, err := OpenFile(input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("read file success!")

	linelist := strings.Split(string(body[:]), "\r\n")

	var outbody string

	for idx, line := range linelist {

		enter := "\r\n"
		if idx+1 == len(linelist) {
			enter = ""
		}

		if len(line) == 0 {
			outbody += fmt.Sprintf(enter)
			continue
		}

		line = swap(line, " ", "")
		line = swap(line, "特使", "Envoy")
		line = swap(line, "侦听器", "监听器")
		line = swap(line, "跨度", "span")

		outbody += fmt.Sprintf("%s%s", addchar(line), enter)
	}

	err = SaveFile(output, []byte(outbody))
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func main() {

	input := "data.txt"
	output := "data2.txt"

	lasttime, _ := os.Stat(input)

	for {

		time.Sleep(1 * time.Second)
		newtime, _ := os.Stat(input)

		if lasttime.ModTime() != newtime.ModTime() {
			process(input, output)
			lasttime, _ = os.Stat(input)
		}
	}

}
