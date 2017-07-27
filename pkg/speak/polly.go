package speak

import (
	"log"
	"github.com/leprosus/golang-tts"
	"os"
	"io/ioutil"
	"fmt"
)

func PollySpeak(text string, options *Options) error {
	log.Printf(">>> Polly: %s", text)

	polly := golang_tts.New(options.Access, options.Secret)
	polly.Format(golang_tts.MP3)

	voice, err := getPollyVoice(options.Language, options.Gender)
	if err != nil {
		return err
	}
	polly.Voice(voice)

	bytes, err := polly.Speech(text)
	if err != nil {
		return err
	}

	mp3, err := ioutil.TempFile("/tmp", "polly")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(mp3.Name(), bytes, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(mp3.Name())

	playCommand := getPlayCommand(mp3.Name())
	_, err = playCommand.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func getPollyVoice(language string, gender string) (string, error) {
	if language == "de" {
		if gender == "female" {
			return golang_tts.Marlene, nil
		}
		return  golang_tts.Hans, nil
	}

	if language == "en" {
		if gender == "female" {
			return golang_tts.Joanna, nil
		}
		return golang_tts.Joey, nil
	}
	return "", fmt.Errorf("Invalid language %s", language)
}
