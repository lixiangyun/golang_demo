package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

func Format(input string) string {

	input = strings.ToLower(input)

	inputbody := []byte(input)

	output := make([]byte, 0)

	for _, v := range inputbody {

		if v == ' ' {
			output = append(output, '-')
		} else if v == '.' {
			continue
		} else if v == '(' || v == ')' {
			continue
		} else {
			output = append(output, v)
		}
	}

	return string(output)
}

func headcheck(b rune) bool {
	if unicode.IsLetter(b) || unicode.IsNumber(b) {
		return true
	}
	return false
}

func main() {

	input := "all.txt"
	output := "all.txt"

	body, err := OpenFile(input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	linelist := strings.Split(string(body[:]), "\r\n")

	var outbody string

	for idx, line := range linelist {

		var benum bool

		if len(line) == 0 {
			if idx+1 != len(linelist) {
				outbody += fmt.Sprintf("\r\n")
			}
			continue
		}

		enum := strings.Index(line, "Enum")
		if enum != -1 {
			benum = true
			line = line[:enum]
		}

		begin := strings.IndexFunc(line, headcheck)
		end := strings.LastIndexFunc(line, headcheck)

		if end == -1 || begin == -1 {
			if idx+1 != len(linelist) {
				outbody += fmt.Sprintf("%s\r\n", line)
			} else {
				outbody += fmt.Sprintf("%s", line)
			}
			continue
		}

		title := line[begin : end+1]

		if benum {
			fmt.Printf("%s (Enum)\r\n", title)
			outbody += fmt.Sprintf("- [%s (Enum)](#%s-enum)\r\n", title, Format(title))
		} else {
			fmt.Printf("%s\r\n", title)
			outbody += fmt.Sprintf("- [%s](#%s)\r\n", title, Format(title))
		}
	}

	err = SaveFile(output, []byte(outbody))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
