package main

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	pid int
	cpu float64
}

func totalcpu() float64 {

	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		processes = append(processes, &Process{pid, cpu})
	}

	var totalcpu float64

	for _, p := range processes {
		//log.Println("Process ", p.pid, " takes ", p.cpu, " % of the CPU")
		totalcpu += p.cpu
	}

	time.Sleep(5 * time.Second)

	return totalcpu
}

func main() {
	var times int
	for {
		cpup := totalcpu()
		if cpup <= 20 {
			times++
		} else {
			times = 0
		}
		log.Printf("total cpu : %.2f , times : %d \r\n", cpup, times)
		if times > 10 {
			return
		}
	}
}
