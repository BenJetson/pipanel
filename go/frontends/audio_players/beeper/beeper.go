package beeper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	pipanel "github.com/BenJetson/pipanel/go"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// SampleRate is the sample rate of the beep/speaker. Defaults to 16 kHz.
// If the sample rate of the beep/speaker is different, change the value to
// resample audio clips appropriately.
var SampleRate beep.SampleRate = 16000

// Beeper implements pipanel.AudioPlayer and plays WAV audio clips from the
// library directory specified. Sound events are expected to omit the .wav file
// extension from the Sound field.
type Beeper struct {
	log         *log.Logger
	libraryPath string
}

// New creates a Beeper instance.
func New(log *log.Logger, libraryPath string) *Beeper {
	// Enforce trailing slash, which makes concatenation with filenames easier.
	if libraryPath[len(libraryPath)-1] != '/' {
		libraryPath += "/"
	}

	// Check to make sure that the directory actually exists; panic otherwise.
	d, err := os.Open(libraryPath)

	if os.IsNotExist(err) {
		err = errors.New("directory referenced by libraryPath does not exist")
	}

	if err != nil {
		panic(err)
	}

	d.Close()

	return &Beeper{log, libraryPath}
}

func validateAudioFilename(fileName string) error {
	// Checks to make sure that only one period exists in the file name.
	// Exists for secutiry purposes to ensure that files outside of the library
	// path cannot be accessed (for example "../not_in_library.wav" is bad).
	if strings.Count(fileName, ".") > 0 {
		return fmt.Errorf("illegal filename '%s' contains periods", fileName)
	}

	return nil
}

// PlaySound handles pipanel sound events.
func (b *Beeper) PlaySound(e pipanel.SoundEvent) error {
	if err := validateAudioFilename(e.Sound); err != nil {
		return nil
	}

	pathToFile := b.libraryPath + e.Sound + ".wav"

	f, err := os.Open(pathToFile)

	if err != nil {
		return nil
	}

	streamer, format, err := wav.Decode(f)

	if err != nil {
		return err
	}

	var streamToPlay beep.Streamer = streamer

	if format.SampleRate != SampleRate {
		streamToPlay = beep.Resample(4, format.SampleRate, SampleRate, streamer)
	}

	speaker.Play(streamToPlay)
	b.log.Printf("Playing sound at '%s'", pathToFile)

	return nil
}
