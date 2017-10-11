package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
)

// 写入value到指定的文件
func writeValues(values []int, outfile string) error {
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println("Failed to create the output file ", outfile)

		return err
	}

	defer file.Close()

	for _, v := range values {
		str := strconv.Itoa(v)
		file.WriteString(str + "\n")
	}

	return nil
}

func main() {

	r := make([]byte, 4)

	values := make([]int, 0)

	for i := 0; i < 1000000; i++ {

		// 构造随机数
		_, err := rand.Read(r)
		if err != nil {
			fmt.Println("rand return failed!", err.Error())
			return
		} else {
			var val uint32 = 0
			for _, v := range r {
				val = uint32(v)
				val = val << 4
			}
			values = append(values, int(val))
		}
	}

	writeValues(values, "unsort.dat")
}
