package logfmt

import (
	"github.com/sirupsen/logrus"

	pipanel "github.com/BenJetson/pipanel/go"
)

// RequestIDFormatter annotates logs with the RequestID field set by the
// server package.
type RequestIDFormatter struct {
	*logrus.TextFormatter
}

const requestIDKey = "requestID"

// Format formats the log entry given using the embedded TextFormatter,
// annotating with the request ID field if available.
func (f *RequestIDFormatter) Format(e *logrus.Entry) ([]byte, error) {
	if e.Context != nil {
		if rID := e.Context.Value(pipanel.RequestIDKey); rID != nil {
			// Create copy of given Entry with added fields.
			annotated := e.WithField(requestIDKey, rID)

			// Copy over remaining fields from parent Entry that were not copied
			// by WithFields. Without this step, this data is lost and will not
			// be logged when a requestID is present.
			annotated.Level = e.Level
			annotated.Caller = e.Caller
			annotated.Message = e.Message
			annotated.Buffer = e.Buffer

			// Set e to point to the annotated copy.
			e = annotated
		}
	}

	return f.TextFormatter.Format(e)
}
