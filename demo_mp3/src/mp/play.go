package mp

import "fmt"

type Player interface {
	Play(source string)
}

func Play(source, mtype string) {
	var p Player

	switch mtype {
	case "MP3":
		p = &Mp3Player{}
	case "WAV":
		p = &WavPlayer{}
	default:
		fmt.Println("Unsupported music type", mtype)
		return
	}

	p.Play(source)
}
