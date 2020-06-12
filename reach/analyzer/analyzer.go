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

	analyzedPaths := make([]reach.AnalyzedPath, len(paths))
	for i, p := range paths {
		traffic := p.TrafficForward()
		predictions := ConnectionPredictions(p)

		analyzedPaths[i] = reach.AnalyzedPath{
			Path:                  p,
			Traffic:               traffic,
			ConnectionPredictions: predictions,
		}

	}
	a.logger.Info("analysis successful", "source", source, "destination", destination)

	return reach.NewAnalysis([]reach.Subject{source, destination}, analyzedPaths), nil
}

// ConnectionPredictions inspects the path to predict the viability of a various kinds of connections made across this network path.
func ConnectionPredictions(path reach.Path) reach.ConnectionPredictionSet {
	result := make(reach.ConnectionPredictionSet)
	traffic := path.TrafficForward()

	if traffic.HasProtocol(reach.ProtocolTCP) {
		result[reach.ProtocolTCP] = ConnectionPredictionTCP(path)
	}

	if traffic.HasProtocol(reach.ProtocolUDP) {
		result[reach.ProtocolUDP] = ConnectionPredictionUDP(path)
	}

	if traffic.HasProtocol(reach.ProtocolICMPv4) {
		result[reach.ProtocolICMPv4] = ConnectionPredictionICMPv4(path)
	}

	if traffic.HasProtocol(reach.ProtocolICMPv6) {
		result[reach.ProtocolICMPv6] = ConnectionPredictionICMPv6(path)
	}

	return result
}

// ConnectionPredictionTCP inspects the path to predict the viability of a TCP connection made across this network path.
func ConnectionPredictionTCP(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionReturnTrafficRequired(path, reach.ProtocolTCP)
}

func connectionPredictionReturnTrafficRequired(path reach.Path, protocol reach.Protocol) reach.ConnectionPrediction {
	failurePossible := false

	for _, segment := range path.Segments() {
		protocolContent := segment.TrafficReturn().Protocol(protocol)
		switch {
		case protocolContent.Empty():
			return reach.ConnectionPredictionFailure
		case protocolContent.Complete() == false:
			failurePossible = true
		}
	}

	if failurePossible {
		return reach.ConnectionPredictionPossibleFailure
	}

	return reach.ConnectionPredictionSuccess
}

func connectionPredictionReturnTrafficOptional(path reach.Path, protocol reach.Protocol) reach.ConnectionPrediction {
	tcs := returnTrafficForSegments(path)
	for _, traffic := range tcs {
		protocolContent := traffic.Protocol(protocol)
		if protocolContent.Complete() == false {
			return reach.ConnectionPredictionPossibleFailure
		}
	}

	return reach.ConnectionPredictionSuccess
}

func returnTrafficForSegments(path reach.Path) []reach.TrafficContent {
	var result []reach.TrafficContent
	for _, segment := range path.Segments() {
		result = append(result, segment.TrafficReturn())
	}
	return result
}

// ConnectionPredictionUDP inspects the path to predict the viability of a UDP connection made across this network path.
func ConnectionPredictionUDP(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionReturnTrafficOptional(path, reach.ProtocolUDP)
}

// ConnectionPredictionICMPv4 inspects the path to predict the viability of an ICMPv4 interaction across this network path.
func ConnectionPredictionICMPv4(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionICMP(path, reach.ProtocolICMPv4)
}

// ConnectionPredictionICMPv6 inspects the path to predict the viability of an ICMPv6 interaction across this network path.
func ConnectionPredictionICMPv6(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionICMP(path, reach.ProtocolICMPv6)
}

func connectionPredictionICMP(path reach.Path, icmpProtocol reach.Protocol) reach.ConnectionPrediction {
	return connectionPredictionReturnTrafficRequired(path, icmpProtocol)
}
