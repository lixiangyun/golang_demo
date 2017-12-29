package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func format(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err.Error())
		return
	}

	buf := make([]byte, 0)

	for {

		var buftmp [128]byte

		cnt, err := file.Read(buftmp[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err.Error())
			return
		}

		buf = append(buf, buftmp[:cnt]...)
	}

	file.Close()

	var output string

	strbuf := strings.Split(string(buf), "\r\n")

	bcode := false

	for _, v := range strbuf {

		if len(v) == 0 {
			output += fmt.Sprintf("\r\n")
			continue
		}

		if bcode {
			output += fmt.Sprintf("%s\r\n", v)
			if v[0] == '}' || v[0] == ']' {
				output += fmt.Sprintf("```\r\n")
				bcode = false
			}
			continue
		}

		switch v[0] {
		case '(':
			{
				tail := strings.Index(v, ")")
				if -1 == tail {
					log.Println("Tail -1")
					return
				}

				header := string(v[1:tail])

				link := "#"

				if -1 != strings.Index(header, "DEFAULT") {
					output += fmt.Sprintf("\t(%s)%s\r\n\r\n", header, string(v[tail+1:]))
				} else {
					cat := strings.Index(header, ",")

					if -1 == cat {
						output += fmt.Sprintf("\t([%s](%s))%s\r\n\r\n", header, link, string(v[tail+1:]))
					} else {
						output += fmt.Sprintf("\t([%s](%s),%s)%s\r\n\r\n", string(header[:cat]), link, string(header[cat+1:]), string(v[tail+1:]))
					}
				}

				bcode = false
			}
		case '[':
			{
				bcode = true
				output += fmt.Sprintf("```\r\n%s\r\n", v)

				if -1 != strings.Index(v, "]") {
					output += fmt.Sprintf("```\r\n")
					bcode = false
				}
			}
		case '{':
			{
				bcode = true
				output += fmt.Sprintf("```\r\n%s\r\n", v)

				if -1 != strings.Index(v, "}") {
					output += fmt.Sprintf("```\r\n")
					bcode = false
				}
			}
		default:
			{
				if v[0] >= 'a' && v[0] <= 'z' && strings.Index(v, ".") == -1 {
					output += fmt.Sprintf("- **%s**<br />\r\n", v)
				} else if strings.Index(v, ".") != -1 && strings.Index(v, " ") != -1 {
					output += fmt.Sprintf("%s\r\n", v)
				} else {
					output += fmt.Sprintf("### %s\r\n", v)
				}

				bcode = false
			}
		}
	}

	fmt.Print(output)

	file, err = os.OpenFile(filename, os.O_WRONLY, 0)
	if err != nil {
		log.Println(err.Error())
		return
	}

	file.WriteString(output)

	file.Sync()

	file.Close()
}

func main() {

	filename := "all.txt"

	format(filename)

	lastinfo, err := os.Stat(filename)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for {

		time.Sleep(1 * time.Second)

		tempinfo, err := os.Stat(filename)
		if err != nil {
			log.Println(err.Error())
			return
		}

		if tempinfo.ModTime() != lastinfo.ModTime() {

			format(filename)

			lastinfo, err = os.Stat(filename)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}
