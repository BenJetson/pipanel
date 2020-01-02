package gtkalerter

import (
	"fmt"
	"time"

	"github.com/gotk3/gotk3/pango"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/BenJetson/humantime"
	pipanel "github.com/BenJetson/pipanel/go"
)

type alertWindow struct {
	window     *gtk.Window
	headerBar  *gtk.HeaderBar
	topLayout  *gtk.Box
	boxLayout  *gtk.Box
	dismissBtn *gtk.Button
	progress   *gtk.ProgressBar
	label      *gtk.Label
	icon       *gtk.Image
	timestamp  time.Time
	inactive   bool
	onDestroy  func()
}

func newAlertWindow(a pipanel.AlertEvent, onDestroy func()) (*alertWindow, error) {
	var w alertWindow
	var err error

	w.timestamp = time.Now()
	w.onDestroy = onDestroy

	// Create the window.
	if w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL); err != nil {
		return nil, err
	}

	w.window.SetPosition(gtk.WIN_POS_CENTER_ALWAYS)
	w.window.SetDecorated(false)
	w.window.SetSizeRequest(450, 300)

	// Create the headerbar.
	if w.headerBar, err = gtk.HeaderBarNew(); err != nil {
		return nil, err
	}

	w.headerBar.SetHasSubtitle(true)
	w.headerBar.SetShowCloseButton(false)
	w.headerBar.SetTitle("Alert")

	// Create the progress bar.
	if w.progress, err = gtk.ProgressBarNew(); err != nil {
		return nil, err
	}

	w.progress.SetFraction(1.0)
	w.progress.SetPulseStep(0.05)

	// Create the message label.
	if w.label, err = gtk.LabelNew("Message"); err != nil {
		return nil, err
	}

	w.label.SetSelectable(false)
	w.label.SetJustify(gtk.JUSTIFY_CENTER)
	w.label.SetLineWrap(true)
	w.label.SetLineWrapMode(pango.WRAP_WORD)

	// Create the icon.
	if w.icon, err = gtk.ImageNewFromIconName(a.Icon, gtk.ICON_SIZE_DIALOG); err != nil {
		return nil, err
	}

	w.icon.SetPixelSize(128)

	// Create the dismiss button.
	if w.dismissBtn, err = gtk.ButtonNewWithLabel("Acknowledge"); err != nil {
		return nil, err
	}

	var ackIcon *gtk.Image
	if ackIcon, err = gtk.ImageNewFromIconName("gtk-yes", gtk.ICON_SIZE_BUTTON); err != nil {
		return nil, err
	}

	w.dismissBtn.SetImage(ackIcon)
	w.dismissBtn.SetAlwaysShowImage(true)

	// Create the layouts.
	if w.boxLayout, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5); err != nil {
		return nil, err
	}

	w.boxLayout.SetMarginStart(15)
	w.boxLayout.SetMarginEnd(15)
	w.boxLayout.SetMarginBottom(15)

	if w.topLayout, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0); err != nil {
		return nil, err
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
	glib.TimeoutAdd(1000, func() bool {
		if w.inactive {
			return false
		}

		w.updateSubtitle()
		return true
	})

	// Fill in the values from the Alert event.
	w.setText(a.Message)
	if !a.Perpetual {
		w.setTimeout(time.Millisecond * a.Timeout)
	} else {
		w.pulseProgress()
	}

	// Register events.
	w.window.Connect("destroy", w.Deactivate)
	w.window.Connect("delete-event", w.Deactivate)
	w.dismissBtn.Connect("clicked", w.Destroy)

	return &w, nil
}

func (w *alertWindow) updateSubtitle() { w.headerBar.SetSubtitle(humantime.Since(w.timestamp)) }

func (w *alertWindow) ShowAll() {
	w.window.ShowAll()
	w.window.SetKeepAbove(true)
}

func (w *alertWindow) Deactivate() { w.inactive = true }

func (w *alertWindow) Destroy() {
	if !w.inactive {
		// Destroy the window.
		w.Deactivate()
		w.window.Destroy()

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

		// Call the onDestroy handler.
		w.onDestroy()
	}
}

func (w *alertWindow) setText(text string) {
	w.label.SetMarkup(fmt.Sprintf(`<span size='36000'>%s</span>`, text))
}

func (w *alertWindow) setTimeout(d time.Duration) {
	expiryTime := w.timestamp.Add(d)

	glib.TimeoutAdd(33, func() bool {
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
}

func (w *alertWindow) pulseProgress() {
	glib.TimeoutAdd(80, func() bool {
		if w.inactive {
			return false
		}

		w.progress.Pulse()

		return true
	})
}
