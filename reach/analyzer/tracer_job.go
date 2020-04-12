package analyzer

import (
	"github.com/luhring/reach/reach"
)

type tracerJob struct {
	partial        reach.PartialPath
	destinationRef reach.InfrastructureReference
}
