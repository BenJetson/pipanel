package gtkalerter

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/pango"
	"github.com/pkg/errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/BenJetson/humantime"
	pipanel "github.com/BenJetson/pipanel/go"
)

type alertWindow struct {
	window       *gtk.Window
	headerBar    *gtk.HeaderBar
	topLayout    *gtk.Box
	boxLayout    *gtk.Box
	dismissBtn   *gtk.Button
	progress     *gtk.ProgressBar
	label        *gtk.Label
	icon         *gtk.Image
	timestamp    time.Time
	inactive     bool
	afterCleanup func()
}

// newAlertWindow creates a new alert window instance. Since Glade is not used
// for layout, this function is long as it must set up each UI element manually.
// nolint: gocyclo
func newAlertWindow(cfg *Config, a pipanel.AlertEvent, afterCleanup func()) (*alertWindow, error) {
	var w alertWindow
	var err error

	w.timestamp = time.Now()
	w.afterCleanup = afterCleanup

	// Create the window.
	if w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk window")
	}

	w.window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
	w.window.SetDecorated(false)
	w.window.SetSizeRequest(cfg.WindowSize.Width, cfg.WindowSize.Height)

	// Create the headerbar.
	if w.headerBar, err = gtk.HeaderBarNew(); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk header bar")
	}

	w.headerBar.SetHasSubtitle(true)
	w.headerBar.SetShowCloseButton(false)
	w.headerBar.SetTitle("Alert")

	// Create the progress bar.
	if w.progress, err = gtk.ProgressBarNew(); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk progress bar")
	}

	w.progress.SetFraction(1.0)
	w.progress.SetPulseStep(0.05)

	// Create the message label.
	if w.label, err = gtk.LabelNew("Message"); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk label for message")
	}

	w.label.SetSelectable(false)
	w.label.SetJustify(gtk.JUSTIFY_CENTER)
	w.label.SetLineWrap(true)
	w.label.SetLineWrapMode(pango.WRAP_WORD)

	// Create the icon.
	if w.icon, err = gtk.ImageNewFromIconName(a.Icon, gtk.ICON_SIZE_DIALOG); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk icon")
	}

	w.icon.SetPixelSize(cfg.IconSize)

	// Create the dismiss button.
	if w.dismissBtn, err = gtk.ButtonNewWithLabel("Acknowledge"); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk button for dismiss")
	}

	var ackIcon *gtk.Image
	if ackIcon, err = gtk.ImageNewFromIconName("gtk-yes", gtk.ICON_SIZE_BUTTON); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk icon for dismiss button")
	}

	w.dismissBtn.SetImage(ackIcon)
	w.dismissBtn.SetAlwaysShowImage(true)

	// Create the layouts.
	if w.boxLayout, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5); err != nil {
		return nil, errors.Wrap(err, "failed to create gtk box layout for alert content")
	}

	w.boxLayout.SetMarginStart(15)
	w.boxLayout.SetMarginEnd(15)
	w.boxLayout.SetMarginBottom(15)

	if w.topLayout, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0); err != nil {
		return nil, errors.Wrap(err, "failed to create box layout for window")
	}

	w.topLayout.SetHomogeneous(false)

	// Add dismiss button to headerbar.
	w.headerBar.PackStart(w.dismissBtn)

	// Add widgets to the box layout.
	w.boxLayout.PackStart(w.icon, false, true, 24)
	w.boxLayout.PackStart(w.progress, false, true, 4)
	w.boxLayout.PackStart(w.label, true, true, 0)

	// Combine layouts and add to the window.
	w.topLayout.Add(w.headerBar)
	w.topLayout.Add(w.boxLayout)
	w.window.Add(w.topLayout)

	// Update the timestamp of this window once per second.
	w.updateSubtitle()
	_, err = glib.TimeoutAdd(1000, func() bool {
		if w.inactive {
			return false
		}

		w.updateSubtitle()
		return true
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add timer for subtitle updates")
	}

	// Fill in the values from the Alert event.
	w.setText(a.Message, cfg.FontSize)
	if !a.Perpetual {
		err = w.setTimeout(time.Millisecond * a.Timeout)
	} else {
		err = w.pulseProgress()
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to add timer for progress bar updates")
	}

	// Register events.
	if _, err = w.window.Connect("destroy", w.Deactivate); err != nil {
		return nil, errors.Wrap(err, "failed to bind destroy signal to window deactivation")
	}
	if _, err = w.window.Connect("delete-event", w.Deactivate); err != nil {
		return nil, errors.Wrap(err, "failed to bind deletion signal to window deactivation")
	}
	if _, err = w.dismissBtn.Connect("clicked", w.Destroy); err != nil {
		return nil, errors.Wrap(err, "failed to bind dismiss button to window destruction")
	}

	return &w, nil
}

func (w *alertWindow) updateSubtitle() { w.headerBar.SetSubtitle(humantime.Since(w.timestamp)) }

func (w *alertWindow) ShowAll() {
	w.window.ShowAll()
	w.window.SetKeepAbove(true)
	w.window.Present()
}

func (w *alertWindow) Deactivate() {
	if !w.inactive {
		w.inactive = true
		w.Cleanup()
	}
}

func (w *alertWindow) Destroy() {
	if !w.inactive {
		// Destroy the window.
		w.window.Destroy()
		w.Deactivate()
	}
}

func (w *alertWindow) Cleanup() {
	// Clear pointers to components so the garbage collector will pick them
	// up and deallocate the (now unreferenced) objects.
	w.window = nil
	w.headerBar = nil
	w.topLayout = nil
	w.boxLayout = nil
	w.dismissBtn = nil
	w.progress = nil
	w.label = nil
	w.icon = nil

	// Call the afterCleanup handler.
	w.afterCleanup()
}

func (w *alertWindow) setText(text string, fontSize int) {
	w.label.SetMarkup(fmt.Sprintf(`<span size='%d000'>%s</span>`,
		fontSize, text))
}

func (w *alertWindow) setTimeout(d time.Duration) error {
	expiryTime := w.timestamp.Add(d)

	_, err := glib.TimeoutAdd(33, func() bool {
		if w.inactive {
			return false
		}

		if time.Now().After(expiryTime) {
			w.Destroy()
			return false
		}

		since := time.Since(w.timestamp)
		frac := 1 - (float64(since) / float64(d))

		w.progress.SetFraction(frac)

		return true
	})

	return errors.Wrap(err, "failed to add timer for updating progress bar fraction")
}

func (w *alertWindow) pulseProgress() error {
	_, err := glib.TimeoutAdd(80, func() bool {
		if w.inactive {
			return false
		}

		w.progress.Pulse()

		return true
	})

	return errors.Wrap(err, "failed to add time for pulsing indeterminate progress bar")
}
