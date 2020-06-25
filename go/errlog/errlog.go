package errlog

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const errKey = "error"
const stackKey = "stack"

// WithError returns a *logrus.Entry derived from the given *logrus.Entry, with
// added fields for the error message and (if present) stack trace.
//
// Caller must ensure that err is not nil.
func WithError(l *logrus.Entry, err error) *logrus.Entry {
	f := logrus.Fields{errKey: err.Error()}

	if stackErr, ok := err.(interface {
		StackTrace() errors.StackTrace
	}); ok {
		f[stackKey] = fmt.Sprintf("%+v", stackErr.StackTrace())
	}

	return l.WithFields(f)
}
