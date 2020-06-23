package gtkalerter

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"

	pipanel "github.com/BenJetson/pipanel/go"
)

// Config specifies the options that modify the behavior of GTKAlerter.
type Config struct {
	// FontSize is the font size that should be used for the alert text content,
	// measured in in points.
	FontSize int `json:"font_size"`
	// WindowSize defines the size of the alert window.
	WindowSize struct {
		// Width defines the width of the alert window, in pixels.
		Width int `json:"width"`
		// Height defines the height of the window, in pixels.
		Height int `json:"height"`
	} `json:"window_size"`
	// IconSize defines the size of the icon, in pixels.
	IconSize int `json:"icon_size"`
	// Defaults contains the values that will replace zero values on the
	// incoming pipanel.AlertEvents.
	Defaults struct {
		// Timeout is the default timeout value, measured in seconds.
		// If a TimeoutRange is specified, this value MUST be within that range.
		Timeout time.Duration `json:"timeout"`
		// Icon is the default icon value, which is the name of a gtk icon.
		Icon string `json:"icon"`
	} `json:"defaults"`
	// TimeoutRange controls the range of values that are acceptable for
	// the timeout field.
	TimeoutRange struct {
		// Min is the minimum timeout value. If this value is zero, the lower
		// bound will not be checked.
		Min time.Duration `json:"min"`
		// Max is the maximum timeout value. If this value is zero, the upper
		// bound will not be checked.
		Max time.Duration `json:"max"`
	} `json:"timeout_range"`
	// ForbidPerpetual determines whether or not perpetual alerts are allowed.
	// If this flag is true, the perpetual flag of incoming alerts will be
	// ignored.
	ForbidPerpetual bool `json:"forbid_perpetual"`
}

// nolint: gocyclo // keeping this validation logic together makes sense
func validateConfig(cfg *Config) error {
	if cfg.FontSize < 1 {
		return errors.New("font size must be set to a value greater than zero")
	}

	if len(cfg.Defaults.Icon) < 1 {
		return errors.New("default icon must be set")
	}

	if cfg.TimeoutRange.Min < 0 {
		return errors.New("min timeout cannot be less than zero")
	} else if cfg.TimeoutRange.Max < 0 {
		return errors.New("max timeout cannot be less than zero")
	} else if cfg.TimeoutRange.Min > cfg.TimeoutRange.Max {
		return errors.New("min timeout cannot be greater than max timeout")
	}

	if cfg.Defaults.Timeout < 1 {
		return errors.New("default timeout cannot be less than one")
	}

	if (cfg.TimeoutRange.Min != 0 && cfg.Defaults.Timeout < cfg.TimeoutRange.Min) ||
		(cfg.TimeoutRange.Max != 0 && cfg.Defaults.Timeout <= cfg.TimeoutRange.Max) {
		return errors.New("default timeout must fall within TimeoutRange")
	}

	if cfg.WindowSize.Width < 250 {
		return errors.New("window cannot be less than 250 pixels wide")
	} else if cfg.WindowSize.Height < 250 {
		return errors.New("window cannot be less than 250 pixels high")
	}

	if cfg.IconSize < 8 {
		return errors.New("icon size cannot be smaller than 8 pixels")
	}

	return nil
}

// GUI is a GTK application that is capable of responding to PiPanel alert
// events by showing them on the screen.
type GUI struct {
	log        *log.Logger
	windowsMux sync.Mutex
	windows    []*alertWindow
	cfg        Config
}

// New creates a fresh GUI instance.
func New() *GUI { return &GUI{} }

func sanitizeAlert(cfg *Config, e *pipanel.AlertEvent) {
	if cfg.ForbidPerpetual {
		e.Perpetual = false
	}

	if e.Timeout == 0 {
		e.Timeout = cfg.Defaults.Timeout
	} else if cfg.TimeoutRange.Min != 0 && e.Timeout < cfg.TimeoutRange.Min {
		e.Timeout = cfg.TimeoutRange.Min
	} else if cfg.TimeoutRange.Max != 0 && e.Timeout > cfg.TimeoutRange.Max {
		e.Timeout = cfg.TimeoutRange.Max
	}

	if len(e.Icon) < 1 {
		e.Icon = cfg.Defaults.Icon
	}
}

// ShowAlert handles alert events by displaying a window to alert the user.
func (g *GUI) ShowAlert(e pipanel.AlertEvent) error {
	sanitizeAlert(&g.cfg, &e)

	_, err := glib.IdleAdd(func() {
		g.log.Println("Waiting for exclusive lock on window list...")

		g.windowsMux.Lock()
		defer g.windowsMux.Unlock()

		g.log.Println("Lock acquired.")

		w, err := newAlertWindow(&g.cfg, e, g.removeInactiveWindows)

		if err != nil {
			err = errors.Wrap(err, "failed to create alert window")
			g.log.Printf("ERROR when creating alert window: %v", err)
			return
		}

		g.log.Println("Displaying alert window to user.")
		w.ShowAll()

		g.windows = append(g.windows, w)
	})

	return errors.Wrap(err, "failed to request creating alert window at next idle")
}

// Init initializes this GUI instance, setting the logger and starting the GTK
// main event loop in a separate goroutine.
func (g *GUI) Init(log *log.Logger, rawCfg json.RawMessage) error {
	g.log = log

	// Decode the config.
	d := json.NewDecoder(bytes.NewReader(rawCfg))
	d.DisallowUnknownFields()

	if err := d.Decode(&g.cfg); err != nil {
		return errors.Wrap(err, "malformed JSON for GTKAlerter configuration")
	}

	if err := validateConfig(&g.cfg); err != nil {
		return errors.Wrap(err, "invalid configuration")
	}

	// Start the GTK main event loop.
	gtk.Init(nil)
	go gtk.Main()

	return nil
}

func (g *GUI) removeInactiveWindows() {
	g.log.Println("Waiting for exclusive lock on window list...")

	g.windowsMux.Lock()
	defer g.windowsMux.Unlock()

	g.log.Println("Lock acquired.")
	g.log.Println("Clearing inactive windows...")

	count := 0

	for i := len(g.windows) - 1; i > -1; i-- {
		if g.windows[i].inactive {
			g.windows = append(g.windows[:i], g.windows[i+1:]...)

			count++
		}
	}

	g.log.Printf("Cleared all inactive windows: %d total.\n", count)
}

// Cleanup tears down this GUI instance, destroying all windows and halting the
// GTK main event loop.
func (g *GUI) Cleanup() error {
	g.log.Println("Shutting down GUI...")

	g.log.Println("Waiting for exclusive lock on window list...")

	g.windowsMux.Lock()
	defer g.windowsMux.Unlock()

	g.log.Println("Lock acquired.")
	g.log.Println("Destroying all windows... ")

	for _, w := range g.windows {
		w.Destroy()
	}

	g.windows = nil

	g.log.Println("All windows destroyed.")

	g.log.Println("Shutting down GTK main event loop...")
	gtk.MainQuit()
	g.log.Println("GTK main event loop halted.")

	return nil
}
