package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang_demo/demo_mp3/src/mlib"
	"golang_demo/demo_mp3/src/mp"
)

var lib *mlib.MusicManager
var id int = 1
var ctrl, signal chan int

func handeleLibCommands(tokens []string) {
	switch tokens[1] {
	case "list":
		{
			for i := 0; i < lib.Len(); i++ {
				e, _ := lib.Get(i)

				fmt.Println(i+1, ":", e.Name, e.Artist, e.Source, e.Type)
			}
		}

	case "add":
		{
			if len(tokens) == 6 {
				id++
				lib.Add(&mlib.MusicEntry{strconv.Itoa(id),
					tokens[2], tokens[3], tokens[4], tokens[5]})
			} else {
				fmt.Println("Usage: lib add <name><artist><source><type>")
			}
		}
	default:
		fmt.Println("Unrecongnized lib command:", tokens[1])
	}
}

func handlePlayCommand(tokens []string) {
	if len(tokens) != 2 {
		fmt.Println("Usage: play <name>")
		return
	}

	e := lib.Find(tokens[1])
	if e == nil {
		fmt.Println("The music ", tokens[1], "does not exist.")
		return
	}

	mp.Play(e.Source, e.Type)
}

func main() {
	fmt.Println(`
		Enter following commands to control the player:
		lib list -- View the existing music lib
		lib add <name><artist><source><type> -- Add a music to the music lib
		lib remove <name> -- Remove the specified music from the lib
		play <name> -- Play the specified music 
		`)

	lib = mlib.NewMusicManager()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command ->")
		rawline, _, _ := r.ReadLine()

		line := string(rawline)

		if line == "q" || line == "e" {
			break
		}

		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "lib":
			handeleLibCommands(tokens)
		case "play":
			handlePlayCommand(tokens)
		default:
			fmt.Println("Unregcognized command:", tokens[0])
		}
	}
}
