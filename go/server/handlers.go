package server

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/logfmt"
)

func parseAndDecodeBody(body io.ReadCloser, target interface{}) error {
	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return errors.Wrap(err, "could not read bytes from request body")
	}

	d := json.NewDecoder(bytes.NewReader(bodyBytes))
	d.DisallowUnknownFields()

	err = d.Decode(target)
	return errors.Wrap(err, "malformed JSON in request body")
}

func (s *Server) handleError(err error, message string, w http.ResponseWriter, statusCode int) bool {
	if err == nil {
		return false
	}

	logfmt.WithError(s.log, err).WithFields(logrus.Fields{
		"code": statusCode,
	}).Errorln("Problem when handling request.")

	http.Error(w, message, statusCode)

	return true
}

func (s *Server) handleAlertEvent(w http.ResponseWriter, r *http.Request) {
	s.log.WithContext(r.Context()).Println("Handling alert event.")

	var e pipanel.AlertEvent
	err := parseAndDecodeBody(r.Body, &e)

	// AlertEvent timeout is measured in milliseconds.
	e.Timeout *= time.Millisecond

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	if !s.processSoundEvent(e.SoundEvent, r, w) {
		return
	}

	err = s.frontend.ShowAlert(r.Context(), e)

	if s.handleError(err, "Failed to present alert to user.", w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSoundEvent(w http.ResponseWriter, r *http.Request) {
	s.log.WithContext(r.Context()).Println("Handling sound event.")

	var e pipanel.SoundEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	if !s.processSoundEvent(e, r, w) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) processSoundEvent(e pipanel.SoundEvent,
	r *http.Request, w http.ResponseWriter) bool {

	if len(e.Sound) < 1 {
		s.log.WithContext(r.Context()).Println("Ignoring empty sound event.")
		return true
	}

	err := s.frontend.PlaySound(r.Context(), e)

	return !s.handleError(err, "Failed to play sound.", w, http.StatusInternalServerError)
}

func (s *Server) handlePowerEvent(w http.ResponseWriter, r *http.Request) {
	s.log.WithContext(r.Context()).Println("Handling power event.")

	var e pipanel.PowerEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	err = s.frontend.DoPowerAction(r.Context(), e)

	if s.handleError(err, "Failed to perform requested power action.", w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleBrightnessEvent(w http.ResponseWriter, r *http.Request) {
	s.log.WithContext(r.Context()).Println("Handling brightness event.")

	var e pipanel.BrightnessEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	err = s.frontend.SetBrightness(r.Context(), e)

	if s.handleError(err, "Failed to perform requested brightness action.", w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}
