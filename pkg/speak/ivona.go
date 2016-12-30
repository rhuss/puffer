package speak

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"runtime"

	ivona "github.com/jpadilla/ivona-go"
)

// Speak converts a text to audio and the send it out via audio
func Speak(text string, options *Options) error {
	client := ivona.New(options.Access, options.Secret)
	r, err := client.CreateSpeech(speechOptions(text, options.Language, options.Gender))
	log.Printf("Received answer: %v", r.ContentType)
	if err != nil {
		log.Fatal(err)
	}

	mp3, err := ioutil.TempFile("/tmp", "ivona")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(mp3.Name(), r.Audio, 0644)
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

func getPlayCommand(mp3 string) *exec.Cmd {
	if runtime.GOOS == "darwin" {
		return exec.Command("afplay", mp3)
	}
	return exec.Command("mpg123", mp3)
}

func speechOptions(text string, language string, gender string) ivona.SpeechOptions {
	voice, err := createVoice(language, gender)
	if err != nil {
		log.Fatal(err)
	}
	return ivona.SpeechOptions{
		Input: &ivona.Input{
			Data: text,
			Type: "text/plain",
		},
		OutputFormat: &ivona.OutputFormat{
			Codec:      "MP3",
			SampleRate: 22050,
		},
		Parameters: &ivona.Parameters{
			Rate:           "medium",
			Volume:         "loud",
			SentenceBreak:  500,
			ParagraphBreak: 640,
		},
		Voice: voice,
	}
}

func createVoice(language string, gender string) (*ivona.Voice, error) {
	if language == "de" {
		if gender == "female" {
			return &ivona.Voice{
				Name:     "Marlene",
				Language: "de-DE",
				Gender:   "Female",
			}, nil
		}
		return &ivona.Voice{
			Name:     "Hans",
			Language: "de-DE",
			Gender:   "Male",
		}, nil
	}

	if language == "en" {
		if gender == "female" {
			return &ivona.Voice{
				Name:     "Amy",
				Language: "en-GB",
				Gender:   "Female",
			}, nil
		}
		return &ivona.Voice{
			Name:     "Brian",
			Language: "en-GB",
			Gender:   "Male",
		}, nil
	}
	return nil, fmt.Errorf("Invalid language %s", language)
}
