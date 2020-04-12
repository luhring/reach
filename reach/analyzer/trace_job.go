package analyzer

import (
	"github.com/luhring/reach/reach"
)

type traceJob struct {
	ref            reach.InfrastructureReference
	path           *reach.Path
	edgeTuple      *reach.IPTuple
	destinationRef reach.InfrastructureReference
}