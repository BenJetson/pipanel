package gui

import (
	"errors"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/BenJetson/humantime"
	pipanel "github.com/BenJetson/pipanel/go"
)

type alertWindow struct {
	builder    *gtk.Builder
	window     *gtk.Window
	headerBar  *gtk.HeaderBar
	dismissBtn *gtk.Button
	progress   *gtk.ProgressBar
	label      *gtk.Label
	icon       *gtk.Image
	timestamp  time.Time
	inactive   bool
}

func newAlertWindow(a pipanel.AlertEvent) (*alertWindow, error) {
	var w alertWindow
	var object glib.IObject
	var err error
	var ok bool

	w.timestamp = time.Now()

	// Read in the view layout from the Glade file.
	// TODO: Fix this so that it is directory agnostic.
	if w.builder, err = gtk.BuilderNewFromFile("glade/alert_window.glade"); err != nil {
		return nil, err
	}

	// Attempt to fetch the window.
	if object, err = w.builder.GetObject("AlertWindow"); err != nil {
		return nil, err
	}

	if w.window, ok = object.(*gtk.Window); !ok {
		return nil, errors.New("object ought to of been a window")
	}

	// Attempt to fetch the headerbar.
	if object, err = w.builder.GetObject("HeaderBar"); err != nil {
		return nil, err
	}

	if w.headerBar, ok = object.(*gtk.HeaderBar); !ok {
		return nil, errors.New("object ought to of been a headerbar")
	}

	// Attempt to fetch the progress bar.
	if object, err = w.builder.GetObject("AlertProgress"); err != nil {
		return nil, err
	}

	if w.progress, ok = object.(*gtk.ProgressBar); !ok {
		return nil, errors.New("object ought to of been a progress bar")
	}

	// Attempt to fetch the label.
	if object, err = w.builder.GetObject("AlertText"); err != nil {
		return nil, err
	}

	if w.label, ok = object.(*gtk.Label); !ok {
		return nil, errors.New("object ought to of been a label")
	}

	// Attempt to fetch the icon.
	if object, err = w.builder.GetObject("AlertIcon"); err != nil {
		return nil, err
	}

	if w.icon, ok = object.(*gtk.Image); !ok {
		return nil, errors.New("object ought to of been an image")
	}

	// Attempt to fetch the dismiss button.
	if object, err = w.builder.GetObject("DismissButton"); err != nil {
		return nil, err
	}

	if w.dismissBtn, ok = object.(*gtk.Button); !ok {
		return nil, errors.New("object ought to of been a button")
	}

	// Update the timestamp of this window once per second.
	glib.TimeoutAdd(1000, func() bool {
		if w.inactive {
			return false
		}

		w.headerBar.SetSubtitle(humantime.Since(w.timestamp))
		return true
	})

	// Fill in the values from the Alert event.
	w.setIcon(a.Icon)
	w.setText(a.Message)
	if !a.Perpetual {
		w.setTimeout(time.Millisecond * a.Timeout)
	}

	// Register events.
	w.window.Connect("destroy", w.Deactivate)
	w.window.Connect("delete-event", w.Deactivate)
	w.dismissBtn.Connect("clicked", w.Destroy)

	return &w, nil
}

func (w *alertWindow) ShowAll() {
	w.window.ShowAll()
	w.window.SetKeepAbove(true)
}

func (w *alertWindow) Deactivate() { w.inactive = true }

func (w *alertWindow) Destroy() {
	if !w.inactive {
		w.Deactivate()
		w.window.Destroy()
	}
}

func (w *alertWindow) setIcon(iconName string) { w.icon.SetFromIconName(iconName, gtk.ICON_SIZE_DIALOG) }

func (w *alertWindow) setText(text string) { w.label.SetText(text) }

func (w *alertWindow) setTimeout(d time.Duration) {
	expiryTime := w.timestamp.Add(d)

	glib.TimeoutAdd(100, func() bool {
		if w.inactive {
			return false
		}

		if time.Now().After(expiryTime) {
			w.Destroy()
			return false
		}

		return true
	})
}
