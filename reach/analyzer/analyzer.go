package analyzer

import (
	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/reachlog"
)

// Analyzer performs Reach's central network traffic analysis.
type Analyzer struct {
	resolver reach.DomainClientResolver
	logger   reachlog.Logger
}

// New creates a new Analyzer.
func New(resolver reach.DomainClientResolver, logger reachlog.Logger) *Analyzer {
	return &Analyzer{
		resolver: resolver,
		logger:   logger,
	}
}

// Analyze performs a full analysis of allowed network traffic among the specified subjects.
func (a *Analyzer) Analyze(source, destination reach.Subject) (*reach.Analysis, error) {
	a.logger.Debug("beginning analysis", "source", source, "destination", destination)

	var tracer reach.Tracer = NewTracer(a.resolver, a.logger)
	paths, err := tracer.Trace(source, destination)
	if err != nil {
		a.logger.Error("analysis failed", "source", source, "destination", destination)
		return nil, err
	}
	a.logger.Info("analysis successful", "source", source, "destination", destination)

	return reach.NewAnalysis([]reach.Subject{source, destination}, paths), nil
}
