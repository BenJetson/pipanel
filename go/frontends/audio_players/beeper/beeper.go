package beeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pipanel "github.com/BenJetson/pipanel/go"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Config is the structure for Beeper configuration.
type Config struct {
	// LibraryPath is the absolute path to the directory storing the audio files
	// to be played by Beeper.
	LibraryPath string `json:"library_path"`
}

// SampleRate is the sample rate of the beep/speaker. Defaults to 16 kHz.
// If the sample rate of the beep/speaker is different, change the value to
// resample audio clips appropriately.
var SampleRate beep.SampleRate = 16000

// Beeper implements pipanel.AudioPlayer and plays WAV audio clips from the
// library directory specified. Sound events are expected to omit the .wav file
// extension from the Sound field.
type Beeper struct {
	log *log.Logger
	cfg Config
}

// New creates a Beeper instance.
func New() *Beeper { return &Beeper{} }

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

	pathToFile := b.cfg.LibraryPath + e.Sound + ".wav"

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

func (b *Beeper) Init(log *log.Logger, rawCfg json.RawMessage) error {
	b.log = log

	// Decode config structure.
	d := json.NewDecoder(bytes.NewReader(rawCfg))
	d.DisallowUnknownFields()

	if err := d.Decode(&(b.cfg)); err != nil {
		return err
	}

	// Make sure library path is set.
	if len(b.cfg.LibraryPath) < 1 {
		return fmt.Errorf("must define an audio library path in config")
	}

	// Enforce trailing slash, which makes concatenation with filenames easier.
	if b.cfg.LibraryPath[len(b.cfg.LibraryPath)-1] != '/' {
		b.cfg.LibraryPath += "/"
	}

	// Check to make sure that the directory actually exists; panic otherwise.
	dir, err := os.Open(b.cfg.LibraryPath)

	if os.IsNotExist(err) {
		return fmt.Errorf("no such directory: %s", b.cfg.LibraryPath)
	} else if err != nil {
		return err
	}

	dir.Close()

	return speaker.Init(SampleRate, SampleRate.N(time.Second/10))
}

func (b *Beeper) Cleanup() error {
	return nil
}
