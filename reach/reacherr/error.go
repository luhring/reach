package reacherr

import (
	"fmt"
	"runtime/debug"
)

type ReachErr struct {
	StackTrace string
	Inner      error
	Message    string
}

func New(err error, messagef string, messageArgs ...interface{}) ReachErr {
	return ReachErr{
		StackTrace: string(debug.Stack()),
		Inner:      err,
		Message:    fmt.Sprintf(messagef, messageArgs...),
	}
}

func (e ReachErr) Error() string {
	return e.Message
}
