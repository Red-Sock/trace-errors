package errors

import (
	"fmt"
	"runtime"
)

func Wrap(innerError error, msg string) error {

	err := Error{
		innerError: innerError,
		msg:        msg,
	}

	if enableTracing {
		runtime.Callers(2, err.trace[:])
	}

	return err
}

func Wrapf(err error, msg string, args ...interface{}) error {
	return Wrap(err, fmt.Sprintf(msg, args...))
}
