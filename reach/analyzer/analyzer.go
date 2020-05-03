package analyzer

import (
	"fmt"

	"github.com/luhring/reach/reach"
)

// Analyzer performs Reach's central network traffic analysis.
type Analyzer struct {
	resolver reach.DomainClientResolver
}

// New creates a new Analyzer.
func New(resolver reach.DomainClientResolver) *Analyzer {
	return &Analyzer{
		resolver: resolver,
	}
}

// Analyze performs a full analysis of allowed network traffic among the specified subjects.
func (a *Analyzer) Analyze(source, destination reach.Subject) (*reach.Analysis, error) {
	var tracer reach.Tracer = NewTracer(a.resolver)
	paths, err := tracer.Trace(source, destination)
	if err != nil {
		return nil, fmt.Errorf("unable to complete trace of network paths: %v", err)
	}

	return reach.NewAnalysis([]reach.Subject{source, destination}, paths), nil
}
