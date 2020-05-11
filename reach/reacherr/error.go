package reacherr

import (
	"fmt"
	"runtime/debug"
	"time"
)

type ReachErr struct {
	Time       time.Time
	StackTrace string
	Inner      error
	Message    string
}

func New(time time.Time, err error, messagef string, messageArgs ...interface{}) ReachErr {
	return ReachErr{
		Inner:      err,
		StackTrace: string(debug.Stack()),
		Time:       time,
		Message:    fmt.Sprintf(messagef, messageArgs...),
	}
}

func (e ReachErr) Error() string {
	return e.Message
}
