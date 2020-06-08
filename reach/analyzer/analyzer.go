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

// ConnectionPredictions inspects the path to predict the viability of a various kinds of connections made across this network path.
func ConnectionPredictions(path reach.Path) (map[reach.Protocol]reach.ConnectionPrediction, error) {
	result := make(map[reach.Protocol]reach.ConnectionPrediction)

	tcpPrediction, err := ConnectionPredictionTCP(path)
	if err != nil {
		return nil, err
	}
	result[reach.ProtocolTCP] = tcpPrediction

	udpPrediction, err := ConnectionPredictionUDP(path)
	if err != nil {
		return nil, err
	}
	result[reach.ProtocolUDP] = udpPrediction

	icmpv4Prediction, err := ConnectionPredictionICMPv4(path)
	if err != nil {
		return nil, err
	}
	result[reach.ProtocolICMPv4] = icmpv4Prediction

	icmpv6Prediction, err := ConnectionPredictionICMPv6(path)
	if err != nil {
		return nil, err
	}
	result[reach.ProtocolICMPv6] = icmpv6Prediction

	return result, nil
}

// ConnectionPredictionTCP inspects the path to predict the viability of a TCP connection made across this network path.
func ConnectionPredictionTCP(path reach.Path) (reach.ConnectionPrediction, error) {
	failurePossible := false

	for _, point := range path.Points {
		returnTraffic, err := reach.NewTrafficContentFromIntersectingMultiple(
			reach.TrafficFromFactors(point.FactorsReturn),
		)
		if err != nil {
			return reach.ConnectionPredictionUnknown, err
		}

		content := returnTraffic.Protocol(reach.ProtocolTCP)
		switch {
		case content.Ports == nil || content.Ports.Empty():
			return reach.ConnectionPredictionFailure, nil
		case content.Ports.Complete() == false:
			failurePossible = true
		}
	}

	if failurePossible {
		return reach.ConnectionPredictionPossibleFailure, nil
	}

	return reach.ConnectionPredictionSuccess, nil
}

// ConnectionPredictionUDP inspects the path to predict the viability of a UDP connection made across this network path.
func ConnectionPredictionUDP(path reach.Path) (reach.ConnectionPrediction, error) {
	for _, point := range path.Points {
		returnTraffic, err := reach.NewTrafficContentFromIntersectingMultiple(
			reach.TrafficFromFactors(point.FactorsReturn),
		)
		if err != nil {
			return reach.ConnectionPredictionUnknown, err
		}

		content := returnTraffic.Protocol(reach.ProtocolUDP)
		if content.Ports == nil || content.Ports.Complete() == false {
			return reach.ConnectionPredictionPossibleFailure, nil
		}
	}

	return reach.ConnectionPredictionSuccess, nil
}

// ConnectionPredictionICMPv4 inspects the path to predict the viability of an ICMPv4 interaction across this network path.
func ConnectionPredictionICMPv4(path reach.Path) (reach.ConnectionPrediction, error) {
	return connectionPredictionICMP(path, reach.ProtocolICMPv4)
}

// ConnectionPredictionICMPv6 inspects the path to predict the viability of an ICMPv6 interaction across this network path.
func ConnectionPredictionICMPv6(path reach.Path) (reach.ConnectionPrediction, error) {
	return connectionPredictionICMP(path, reach.ProtocolICMPv6)
}

func connectionPredictionICMP(path reach.Path, icmpProtocol reach.Protocol) (reach.ConnectionPrediction, error) {
	failurePossible := false

	for _, point := range path.Points {
		returnTraffic, err := reach.NewTrafficContentFromIntersectingMultiple(
			reach.TrafficFromFactors(point.FactorsReturn),
		)
		if err != nil {
			return reach.ConnectionPredictionUnknown, err
		}

		content := returnTraffic.Protocol(icmpProtocol)
		switch {
		case content.ICMP == nil || content.ICMP.Empty():
			return reach.ConnectionPredictionFailure, nil
		case content.ICMP.Complete() == false:
			failurePossible = true
		}
	}

	if failurePossible {
		return reach.ConnectionPredictionPossibleFailure, nil
	}

	return reach.ConnectionPredictionSuccess, nil
}
