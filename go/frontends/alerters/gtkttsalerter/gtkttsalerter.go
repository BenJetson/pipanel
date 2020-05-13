package gtkttsalerter

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/frontends/alerters/gtkalerter"
	"github.com/BenJetson/pipanel/go/frontends/alerters/ttsalerter"
)

// Config specifies the options that modify the behavior of GTKAlerter,
// TTSAlerter, and GTKTTSAlerter.
type Config struct {
	TTSAlerterCfg json.RawMessage `json:"tts_alerter"`
	GTKAlerterCfg json.RawMessage `json:"gtk_alerter"`
	// NoTTSPrefix may be prepended to the text of a pipanel.AlertEvent.Message
	// to prevent the TTSAlerter from reading the message out loud.
	//
	// If not set, the TTSAlerter will always read the message out loud.
	// Defaults to not set, so messages are always read out loud by default.
	NoTTSPrefix string `json:"no_tts_prefix"`
}

// GTKTTSAlerter handles PiPanel alert events by displaying them on-screen and
// reading them out loud.
type GTKTTSAlerter struct {
	*gtkalerter.GUI
	*ttsalerter.TTSAlerter
	log         *log.Logger
	cfg         Config
	checkPrefix bool
}

func New() *GTKTTSAlerter {
	return &GTKTTSAlerter{
		GUI:        gtkalerter.New(),
		TTSAlerter: ttsalerter.New(),
	}
}

// Init initializes this GTKTTSAlerter, parsing the configuration and 
// initializing both GTKAlerter and TTSAlerter.
func (g *GTKTTSAlerter) Init(log *log.Logger, rawCfg json.RawMessage) error {
	g.log = log

	// Load config so it can be separated.
	d := json.NewDecoder(bytes.NewReader(rawCfg))
	d.DisallowUnknownFields()

	if err := d.Decode(&(g.cfg)); err != nil {
		return err
	}

	// Enable the checkPrefix feature if a prefix is specified.
	g.checkPrefix = len(g.cfg.NoTTSPrefix) > 0

	// Initialize GTKAlerter and TTSAlerter with their respective configs.
	if err := g.GUI.Init(log, g.cfg.GTKAlerterCfg); err != nil {
		return err
	}

	if err := g.TTSAlerter.Init(log, g.cfg.TTSAlerterCfg); err != nil {
		return err
	}

	return nil
}

// Cleanup tears down this GTKTTSAlerter instance, triggering cleanup of
// the GTKAlerter and TTSAlerter.
func (g *GTKTTSAlerter) Cleanup() error {
	if err := g.GUI.Cleanup(); err != nil {
		return err
	}

	if err := g.TTSAlerter.Cleanup(); err != nil {
		return err
	}

	return nil
}

// ShowAlert displays the alert on the screen using gtkalerter.GUI and
// (provided No TTS prefix is not present) reads the message out loud.
func (g *GTKTTSAlerter) ShowAlert(e pipanel.AlertEvent) error {
	// If a message has the no TTS prefix, it should not be read out loud.
	shouldReadMsg := true
	if g.checkPrefix && strings.HasPrefix(e.Message, g.cfg.NoTTSPrefix) {
		shouldReadMsg = false
		e.Message = e.Message[len(g.cfg.NoTTSPrefix):]
		g.log.Println("Detected No TTS prefix; skipping alert read-out.")
	}

	err := g.GUI.ShowAlert(e)

	if err != nil {
		return err
	}

	if shouldReadMsg {
		return g.TTSAlerter.ShowAlert(e)
	}
	return nil
}
