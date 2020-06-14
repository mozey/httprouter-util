package logutil

// Code below mostly copied from
// https://github.com/rs/zerolog/blob/master/pkgerrors/stacktrace.go

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	return false
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}

// TruncFunc can be used to truncate the stack trace from the specified func.
// Use it to remove functions for setup from the trace,
// this makes the trace easier to read
var TruncFunc = "HandlerFunc.ServeHTTP"

// MarshalStack implements pkg/errors stack trace marshaling.
// The stack is only traced up to but not including TruncFunc.
// Copied from pkgerrors.MarshalStack
func MarshalStack(err error) interface{} {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	sterr, ok := err.(stackTracer)
	if !ok {
		return nil
	}
	st := sterr.StackTrace()
	s := &state{}
	out := make([]string, 0, len(st))
	for _, frame := range st {
		filePaths := strings.Split(fmt.Sprintf("%+s", frame), "\n\t")
		var filePath string
		if len(filePaths) == 2 {
			filePath = filePaths[1]
		} else {
			filePath = frameField(frame, s, 's')
		}
		functionName := frameField(frame, s, 'n')
		lineNumber := frameField(frame, s, 'd')
		if functionName == TruncFunc {
			break
		}
		// Space in front of filePath required for click through to definition
		out = append(out, fmt.Sprintf(" %s:%s %s",
			filePath,
			lineNumber,
			functionName))
	}
	return out
}
