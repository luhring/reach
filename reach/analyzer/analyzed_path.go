package analyzer

import "github.com/luhring/reach/reach"

type AnalyzedPath struct {
	path           *reach.Path
	forwardTraffic reach.TrafficContent
	returnTraffic  reach.TrafficContent
}
