package logfmt

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ErrKey = "error"
const StackKey = "stack"

// WithError returns a *logrus.Entry derived from the given *logrus.Entry, with
// added fields for the error message and (if present) stack trace.
//
// Caller must ensure that err is not nil.
func WithError(l *logrus.Entry, err error) *logrus.Entry {
	f := logrus.Fields{ErrKey: err.Error()}

	if stackErr, ok := err.(interface {
		StackTrace() errors.StackTrace
	}); ok {
		f[StackKey] = fmt.Sprintf("%+v", stackErr.StackTrace())
	}

	return l.WithFields(f)
}
