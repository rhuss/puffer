package speak

import (
	"log"
	"runtime"
	"os/exec"
)

// Speak converts a text to audio and the send it out via audio
func Speak(text string, options *Options) error {
	if options.Backend == "ivona" {
		return IvonaSpeak(text, options)
	}
	if options.Backend == "polly" {
		return PollySpeak(text, options)
	}

	log.Printf(">>> Unknown backend %s. Ignoring", options.Backend)
	return nil
}


func getPlayCommand(mp3 string) *exec.Cmd {
	if runtime.GOOS == "darwin" {
		return exec.Command("afplay", mp3)
	}
	return exec.Command("mpg123", mp3)
}
