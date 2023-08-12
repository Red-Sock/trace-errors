package errors

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

var enableTracing = true

func New(msg string) error {
	err := Error{
		msg: msg,
	}

	if enableTracing {
		runtime.Callers(2, err.trace[:])
	}

	return err
}

type Error struct {
	innerError error
	msg        string
	trace      [3]uintptr
}

func (e Error) Error() string {

	msg := ""
	if e.innerError != nil {
		msg = e.innerError.Error()
	}

	msg += "\n" + e.msg + "\n"

	if enableTracing {
		frames := runtime.CallersFrames(e.trace[:])
		fr, ok := frames.Next()
		if ok {
			msg += "\n" + strings.Join([]string{fr.Function + "()", "        " + fr.File + ":" + strconv.Itoa(fr.Line)}, "\n")
		}
	}

	return msg
}

func Is(err1, err2 error) bool {
	return errors.Is(err1, err2)
}
