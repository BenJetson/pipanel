package server

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	pipanel "github.com/BenJetson/pipanel/go"
)

func parseAndDecodeBody(body io.ReadCloser, target interface{}) error {
	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	d := json.NewDecoder(bytes.NewReader(bodyBytes))
	d.DisallowUnknownFields()
	return d.Decode(target)
}

func (s *Server) handleError(err error, message string, w http.ResponseWriter, statusCode int) bool {
	if err == nil {
		return false
	}

	s.log.Printf("ERROR: (STATUS %d) %s\n", statusCode, err)
	http.Error(w, message, statusCode)

	return true
}

func (s *Server) handleAlertEvent(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Handling alert event.")

	var e pipanel.AlertEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	err = s.frontend.ShowAlert(e)

	if s.handleError(err, "Failed to present alert to user.", w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleSoundEvent(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Handling sound event.")

	var e pipanel.SoundEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	err = s.frontend.PlaySound(e)

	if s.handleError(err, "Failed to play sound.", w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handlePowerEvent(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Handling power event.")

	var e pipanel.PowerEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	err = s.frontend.DoPowerAction(e)

	if s.handleError(err, "Failed to perform requested power action.", w, http.StatusBadRequest) {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleBrightnessEvent(w http.ResponseWriter, r *http.Request) {
	s.log.Println("Handling brightness event.")

	var e pipanel.BrightnessEvent
	err := parseAndDecodeBody(r.Body, &e)

	if s.handleError(err, "JSON is invalid or violates schema.", w, http.StatusBadRequest) {
		return
	}

	w.WriteHeader(http.StatusOK)
}
