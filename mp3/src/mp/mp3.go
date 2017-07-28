package mp

import (
	"fmt"
	"time"
)

type Mp3Player struct {
	stat     int
	progress int
}

func (p *Mp3Player) Play(source string) {
	fmt.Println("Playing MP3 music", source)

	p.progress = 0

	for p.progress < 100 {
		time.Sleep(100 * time.Millisecond)
		fmt.Print(".")
		p.progress += 10
	}

	fmt.Println("\nFinished playing", source)
}
