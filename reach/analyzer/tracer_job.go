package analyzer

import (
	"github.com/luhring/reach/reach"
)

type tracerJob struct {
	source      reach.Subject
	destination reach.Subject
	partial     reach.PartialPath
}
