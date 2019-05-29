package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/* ps aux

USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         2  0.0  0.0      0     0 ?        S    May25   0:00 [kthreadd]
root         3  0.0  0.0      0     0 ?        S    May25   0:14 [ksoftirqd/0]
...

*/

type Process struct {
	user  string
	pid   int
	cpup  float64
	memp  float64
	vsz   int
	rss   int
	tty   string
	stat  string
	start string
	time  string
	name  string

	samples int
}

func (proc *Process) Add(proc2 Process) {
	proc.cpup += proc2.cpup
	proc.memp += proc2.memp
	proc.vsz += proc2.vsz
	proc.rss += proc2.rss
	proc.samples++
}

func (proc *Process) Stat() {
	proc.samples++
	proc.cpup = proc.cpup / float64(proc.samples)
	proc.memp = proc.memp / float64(proc.samples)
	proc.vsz = proc.vsz / (proc.samples)
	proc.rss = proc.rss / (proc.samples)
}

func IgnoreSpaces(line string) string {
	var flag bool
	for i := 0; i < len(line); i++ {
		if line[i] == ' ' || line[i] == '\t' {
			flag = true
		}
		if line[i] != ' ' && line[i] != '\t' && flag {
			return line[i:]
		}
	}
	return ""
}

func CatSpaces(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == ' ' || line[i] == '\t' {
			return line[:i]
		}
	}
	return ""
}

func CatEnd(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == '\n' {
			return line[:i]
		}
	}
	return ""
}

func ParseLine(line string) Process {

	var proc Process

	proc.user = CatSpaces(line)
	line = IgnoreSpaces(line)

	pid := CatSpaces(line)
	proc.pid, _ = strconv.Atoi(pid)
	line = IgnoreSpaces(line)

	cpup := CatSpaces(line)
	proc.cpup, _ = strconv.ParseFloat(cpup, 64)
	line = IgnoreSpaces(line)

	memp := CatSpaces(line)
	proc.memp, _ = strconv.ParseFloat(memp, 64)
	line = IgnoreSpaces(line)

	vsz := CatSpaces(line)
	proc.vsz, _ = strconv.Atoi(vsz)
	line = IgnoreSpaces(line)

	rss := CatSpaces(line)
	proc.rss, _ = strconv.Atoi(rss)
	line = IgnoreSpaces(line)

	proc.tty = CatSpaces(line)
	line = IgnoreSpaces(line)

	proc.stat = CatSpaces(line)
	line = IgnoreSpaces(line)

	proc.start = CatSpaces(line)
	line = IgnoreSpaces(line)

	proc.time = CatSpaces(line)
	line = IgnoreSpaces(line)

	proc.name = CatEnd(line)

	return proc
}

func collectInfo() []Process {
	var out bytes.Buffer

	cmd := exec.Command("ps", "aux")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	procs := make([]Process, 0)

	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		if -1 != strings.Index(line, "COMMAND") && -1 != strings.Index(line, "TIME") {
			continue
		}

		proc := ParseLine(line)
		procs = append(procs, proc)
	}

	return procs
}

func sampleInfo(samples map[int]*Process, proc2 Process) {
	var proc *Process
	proc, b := samples[proc2.pid]
	if !b {
		samples[proc2.pid] = &proc2
	} else {
		proc.Add(proc2)
	}
}

func main() {
	matchs := os.Args[1:]
	samples := make(map[int]*Process, 0)

	for i := 0; i < 100; i++ {
		time.Sleep(5 * time.Second)

		procs := collectInfo()
		for _, proc := range procs {
			if len(matchs) > 0 {
				for _, match := range matchs {
					if -1 != strings.Index(proc.name, match) {
						sampleInfo(samples, proc)
						break
					}
				}
			} else {
				sampleInfo(samples, proc)
			}
		}
		fmt.Printf(".")
	}
	fmt.Println()
	for _, proc := range samples {
		proc.Stat()
		fmt.Printf("%d\t%s\t%d\t%.2f\t%.2f\t%d\t%d\t%s\n",
			proc.samples, proc.user, proc.pid, proc.cpup, proc.memp, proc.vsz, proc.rss, proc.name)
	}
}
