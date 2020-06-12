package reach

import "github.com/luhring/reach/reach/traffic"

// A Factor describes how a particular component of the ingested resources has an impact on the network traffic allowed to flow from a source to a destination.
type Factor struct {
	Kind       string
	Resource   Reference
	Traffic    traffic.Content
	Properties interface{} `json:"Properties,omitempty"`
}

// TrafficFromFactors returns a set of TrafficContents found among the input factors.
func TrafficFromFactors(factors []Factor) []traffic.Content {
	var result []traffic.Content

	for _, factor := range factors {
		result = append(result, factor.Traffic)
	}

	return result
}
