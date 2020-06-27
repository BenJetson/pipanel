package server

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
	"github.com/BenJetson/pipanel/go/logfmt"
)

// AttachRequestIDMiddlewareBuilder attaches a unique request identifier to
// each request via its context.
func AttachRequestIDMiddlewareBuilder() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Attach a new UUID as the request ID.
			r = r.WithContext(context.WithValue(
				r.Context(),
				pipanel.RequestIDKey,
				uuid.New().String(),
			))

			// Continue handling request.
			h(w, r)
		}
	}
}

// PanicRecoverMiddlewareBuilder creates a new Middleware that recovers from
// panics by logging the stack and cause, then returning HTTP 500 to the client.
func PanicRecoverMiddlewareBuilder(log *logrus.Entry) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Send useful debuggin information to the console.
					log.WithContext(r.Context()).WithFields(logrus.Fields{
						logfmt.StackKey: debug.Stack(),
						"cause":         rec,
					}).Errorln("Request handler panicked! Recovering.")

					// Clear the ResponseWriter and send internal error code.
					if wf, ok := w.(http.Flusher); ok {
						wf.Flush()
					}
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			h(w, r)
		}
	}
}
