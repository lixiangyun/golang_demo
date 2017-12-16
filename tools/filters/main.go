package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

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

			strtmp := fmt.Sprintf("(../../%s.md)\r\n", string(output))
			output = []byte(strtmp)

		} else {
			output = append(output, "\r\n"...)
		}

		os.Stdout.Write(output[:])

	}
}
