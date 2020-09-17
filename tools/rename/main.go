package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	HELP bool

	DIR    string
	CFG    string
	Separator string
)

type Config struct {
	Str []string `json:"strings"`
}

func  init()  {
	Separator = fmt.Sprintf("%c",os.PathSeparator)

	flag.StringVar(&CFG, "config", "config.json", "config filename.")
	flag.BoolVar(&HELP,"help",false,"This help.")
	flag.StringVar(&DIR,"dir","."+Separator,"directory path file rename.")
}

func (cfg *Config)RenameFile(filename string)  {
	for _, v := range cfg.Str {
		if len(v) == 0 {
			continue
		}
		idx := strings.Index(filename, v)
		if idx == -1 {
			continue
		}
		newfilename := filename[:idx] + filename[idx+len(v):]
		err := os.Rename(filename, newfilename)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("change file name [%s] to [%s]\n", filename, newfilename)
		}
	}
}

func (cfg *Config)FindFile(dir string) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			cfg.FindFile(dir + Separator + file.Name())
		} else {
			filename := file.Name()
			cfg.RenameFile(dir + Separator + filename)
		}
	}
}

func main() {
	flag.Parse()
	if HELP {
		flag.Usage()
		os.Exit(-1)
	}

	value, err := ioutil.ReadFile(CFG)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	cfg := new(Config)
	err = json.Unmarshal(value, cfg)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if len(cfg.Str) == 0 {
		fmt.Printf("no any to rename")
		os.Exit(-1)
	}

	cfg.FindFile(DIR)
}
