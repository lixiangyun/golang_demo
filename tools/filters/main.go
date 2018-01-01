package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func getprefix() string {

	os.Stdout.Write([]byte("Please input prefix length [0~10]. \r\n"))

	var input [1024]byte
	cnt, err := os.Stdin.Read(input[:])
	if err != nil {
		log.Println("error: ", err.Error())
		return "../"
	}

	num, err := strconv.Atoi(string(input[:cnt-2]))
	if err != nil {
		log.Println("error: ", err.Error())
		return "../"
	}

	var prefix string

	for i := 0; i < num; i++ {
		prefix += "../"
	}

	return prefix
}

func main() {

	prefix := getprefix()

	fmt.Println("use the prefix : ", prefix)

	for {
		var input [1024]byte
		n, err := os.Stdin.Read(input[:])

		if err != nil {
			log.Println("error: ", err.Error())
			break
		}

		output := make([]byte, 0)

		for _, v := range input[:n] {

			var flag bool

			switch v {
			case ' ':
			case '(':
			case ')':
			case '.':
			case '-':
			case '_':
			case '/':
			case 194:
			case '?':
				{
					flag = true
					v = '/'
				}
			case 187:
				{
					flag = true
					v = '/'
				}
			case '\n':
			case '\t':
			case '\r':

			default:
				flag = true
			}

			//log.Printf("%c:%d", v, v)

			if flag {
				output = append(output, v)
			}
		}

		if len(output) != 0 {

			strtmp := fmt.Sprintf("(%s%s.md)\r\n", prefix, string(output))
			output = []byte(strtmp)

		} else {
			output = append(output, "\r\n"...)
		}

		os.Stdout.Write(output[:])

	}
}
