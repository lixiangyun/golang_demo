package mp

import (
	"fmt"
	"time"
)

type WavPlayer struct {
	stat     int
	progress int
}

func (p *WavPlayer) Play(source string) {
	fmt.Println("Playing WAV music", source)

	p.progress = 0

	for p.progress < 100 {
		time.Sleep(100 * time.Millisecond)
		fmt.Print(".")
		p.progress += 10
	}

	fmt.Println("\nFinished playing", source)
}
