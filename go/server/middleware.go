package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	pipanel "github.com/BenJetson/pipanel/go"
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
