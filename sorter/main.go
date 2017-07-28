// sorter project main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang_demo/sorter/src/bubblesort"
	"golang_demo/sorter/src/qsort"
	"io"
	"os"
	"strconv"
	"time"
)

var infile *string = flag.String("i", "unsort.dat", "File contains values for sorting.")
var outfile *string = flag.String("o", "sorted.dat", "File to receive sorted values.")
var algorithm *string = flag.String("a", "qsort", "Sort algorithm.")

func readValues(infile string) (values []int, err error) {

	file, err := os.Open(infile)

	if err != nil {
		fmt.Println("Failed to open the file ", infile)
		return
	}

	defer file.Close()

	br := bufio.NewReader(file)

	values = make([]int, 0, 100)

	for {
		line, isPrefix, err1 := br.ReadLine()

		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}

		if isPrefix {
			fmt.Println("A too long line , seems unexpected.")
			return
		}

		str := string(line)

		if len(str) == 0 {
			continue
		}

		value, err2 := strconv.Atoi(str)

		if err2 != nil {
			err = err2

			return
		}

		values = append(values, value)
	}

	return
}

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

func main1() {

	//flag.Parse()

	values, err := readValues(*infile)

	if err == nil {
		fmt.Println("Read values : ", values)
	} else {
		fmt.Println(err)
	}

	t1 := time.Now()

	values = qsort.QuickSort(values)
	//values = bubblesort.BubbleSort(values)

	t2 := time.Now()

	err1 := writeValues(values, *outfile)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("Write values : ", values)
	}

	fmt.Println("Time Used : ", t2.Sub(t1))
}

func main2() {

	//flag.Parse()

	values, err := readValues(*infile)

	if err == nil {
		fmt.Println("Read values : ", values)
	} else {
		fmt.Println(err)
	}

	//values = qsort.QuickSort(values)

	t1 := time.Now()

	values = bubblesort.BubbleSort(values)

	t2 := time.Now()

	err1 := writeValues(values, *outfile)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("Write values : ", values)
	}

	fmt.Println("Time Used : ", t2.Sub(t1))

}

func main() {
	main1()
	main2()
}
