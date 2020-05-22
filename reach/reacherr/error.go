package reacherr

import (
	"fmt"
	"runtime/debug"
)

// ReachErr is the interface used to wrap known errors within the Reach library.
//
// Consumers of the Reach library should expect that if a ReachErr is returned (determined via an interface check), the error is a known edge case. Any errors returned that are not wrapped in a ReachErr can be considered bugs and should be reported.
type ReachErr interface {
	error
	Inner() error
	StackTrace() string
	Message() string
	Unwrap() []error
}

// New returns a new instance of a ReachErr.
func New(err error, messagef string, messageArgs ...interface{}) ReachErr {
	return &reachErr{
		inner:      err,
		stackTrace: string(debug.Stack()),
		message:    fmt.Sprintf(messagef, messageArgs...),
	}
}

// SafeWrap wraps the provided err with a new ReachErr only if the err is already a ReachErr. If err is not a ReachErr, the err is returned as is.
func SafeWrap(err error, messagef string, messageArgs ...interface{}) error {
	if reachErr, ok := err.(ReachErr); ok {
		wrapper := New(reachErr, messagef, messageArgs...)
		return wrapper
	}
	return err
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

func (e reachErr) Unwrap() []error {
	var errs []error
	errs = append(errs, e)
	errs = append(errs, getInners(e)...)
	return errs
}

func getInners(err ReachErr) []error {
	var errs []error
	if inner := err.Inner(); inner != nil {
		errs = append(errs, inner)
		if innerReachErr, ok := inner.(ReachErr); ok {
			errs = append(errs, getInners(innerReachErr)...)
		}
	}
	return errs
}
