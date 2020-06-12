package analyzer

import (
	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/reachlog"
	"github.com/luhring/reach/reach/traffic"
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
		analyzedPaths[i] = reach.AnalyzedPath{
			Path:                  p,
			Traffic:               p.TrafficForward(),
			ConnectionPredictions: ConnectionPredictions(p),
		}

	}
	a.logger.Info("analysis successful", "source", source, "destination", destination)

	return reach.NewAnalysis([]reach.Subject{source, destination}, analyzedPaths), nil
}

// ConnectionPredictions inspects the path to predict the viability of a various kinds of connections made across this network path.
func ConnectionPredictions(path reach.Path) reach.ConnectionPredictionSet {
	result := make(reach.ConnectionPredictionSet)
	t := path.TrafficForward()

	if t.HasProtocol(traffic.ProtocolTCP) {
		result[traffic.ProtocolTCP] = ConnectionPredictionTCP(path)
	}

	if t.HasProtocol(traffic.ProtocolUDP) {
		result[traffic.ProtocolUDP] = ConnectionPredictionUDP(path)
	}

	if t.HasProtocol(traffic.ProtocolICMPv4) {
		result[traffic.ProtocolICMPv4] = ConnectionPredictionICMPv4(path)
	}

	if t.HasProtocol(traffic.ProtocolICMPv6) {
		result[traffic.ProtocolICMPv6] = ConnectionPredictionICMPv6(path)
	}

	return result
}

// ConnectionPredictionTCP inspects the path to predict the viability of a TCP connection made across this network path.
func ConnectionPredictionTCP(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionReturnTrafficRequired(path, traffic.ProtocolTCP)
}

func connectionPredictionReturnTrafficRequired(path reach.Path, protocol traffic.Protocol) reach.ConnectionPrediction {
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

func connectionPredictionReturnTrafficOptional(path reach.Path, protocol traffic.Protocol) reach.ConnectionPrediction {
	for _, content := range returnTrafficForSegments(path) {
		protocolContent := content.Protocol(protocol)
		if protocolContent.Complete() == false {
			return reach.ConnectionPredictionPossibleFailure
		}
	}

	return reach.ConnectionPredictionSuccess
}

func returnTrafficForSegments(path reach.Path) []traffic.Content {
	var result []traffic.Content
	for _, segment := range path.Segments() {
		result = append(result, segment.TrafficReturn())
	}
	return result
}

// ConnectionPredictionUDP inspects the path to predict the viability of a UDP connection made across this network path.
func ConnectionPredictionUDP(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionReturnTrafficOptional(path, traffic.ProtocolUDP)
}

// ConnectionPredictionICMPv4 inspects the path to predict the viability of an ICMPv4 interaction across this network path.
func ConnectionPredictionICMPv4(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionICMP(path, traffic.ProtocolICMPv4)
}

// ConnectionPredictionICMPv6 inspects the path to predict the viability of an ICMPv6 interaction across this network path.
func ConnectionPredictionICMPv6(path reach.Path) reach.ConnectionPrediction {
	return connectionPredictionICMP(path, traffic.ProtocolICMPv6)
}

func connectionPredictionICMP(path reach.Path, icmpProtocol traffic.Protocol) reach.ConnectionPrediction {
	return connectionPredictionReturnTrafficRequired(path, icmpProtocol)
}
