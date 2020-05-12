package reacherr

import (
	"fmt"
	"runtime/debug"
)

type ReachErr interface {
	error
	Inner() error
	StackTrace() string
	Message() string
}

func New(err error, messagef string, messageArgs ...interface{}) ReachErr {
	return &reachErr{
		inner:      err,
		stackTrace: string(debug.Stack()),
		message:    fmt.Sprintf(messagef, messageArgs...),
	}
}

type reachErr struct {
	inner      error
	stackTrace string
	message    string
}

func (e reachErr) Inner() error {
	return e.inner
}

func (e reachErr) StackTrace() string {
	return e.stackTrace
}

func (e reachErr) Message() string {
	return e.message
}

func (e reachErr) Error() string {
	return e.Message()
}
