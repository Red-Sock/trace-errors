package errors

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var enableTracing = false

var enableTracingFlag = "--enable-rscli-tracing"

func init() {
	// in case when project was compiled with
	// "rscliErrorTracingDisabled" build flag,
	// but we need traces
	for _, item := range os.Args {
		if item == enableTracingFlag {
			enableTracing = true
			return
		}
	}
}

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

func (e Error) Error() (msg string) {
	msg += e.msg

	if e.innerError != nil {

		msg += "\n" + e.innerError.Error()
	}

	if enableTracing {
		frames := runtime.CallersFrames(e.trace[:])
		fr, ok := frames.Next()
		if ok {
			msg += "\n" + strings.Join([]string{fr.Function + "()", "        " + fr.File + ":" + strconv.Itoa(fr.Line)}, "\n")
		}
	}

	return msg
}

func (e Error) Unwrap() error {
	return e.innerError
}

func Is(err1, err2 error) bool {
	return errors.Is(err1, err2)
}
