package reach

// A Factor describes how a particular component of the ingested resources has an impact on the network traffic allowed to flow from a source to a destination.
type Factor struct {
	Kind          string
	Resource      ResourceReference
	Traffic       TrafficContent
	ReturnTraffic TrafficContent // DEPRECATED
	Properties    interface{}    `json:"Properties,omitempty"`
}

func TrafficFromFactors(factors []Factor) []TrafficContent {
	var result []TrafficContent

	for _, factor := range factors {
		result = append(result, factor.Traffic)
	}

	return result
}
