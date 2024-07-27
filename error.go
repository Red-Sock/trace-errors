package errors

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func New(msg string, args ...any) error {
	str, grpcCode := split(args)

	err := Error{
		msg:      strings.Join(append([]string{msg}, str...), "; "),
		grpcCode: grpcCode,
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

	grpcCode *codes.Code
}

func (e Error) Error() (msg string) {
	msg += e.msg

	if e.innerError != nil {
		msg = e.innerError.Error() + "\n" + msg
	}

	if enableTracing {
		frames := runtime.CallersFrames(e.trace[:])
		fr, ok := frames.Next()
		if ok {
			traceStr := strings.Join(
				[]string{fr.Function + "()",
					"        " + fr.File + ":" + strconv.Itoa(fr.Line)}, "\n")
			msg = traceStr + "\n" + msg
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

func (e Error) GRPCStatus() *status.Status {
	if e.grpcCode != nil {
		return status.New(*e.grpcCode, e.Error())
	}

	ie, ok := e.innerError.(Error)
	if ok {
		return ie.GRPCStatus()
	}

	return status.New(codes.Internal, e.Error())
}
