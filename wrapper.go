package errors

import (
	"fmt"
	"runtime"
	"strings"
)

func Wrap(innerError error, msg ...string) error {

	err := Error{
		innerError: innerError,
		msg:        strings.Join(msg, "; "),
	}

	if enableTracing {
		runtime.Callers(2, err.trace[:])
	}

	return err
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(msg, args...))
}
