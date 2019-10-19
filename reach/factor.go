package reach

// A Factor describes how a particular component of the ingested resources has an impact on the network traffic allowed to flow from a source to a destination.
type Factor struct {
	Kind       string
	Resource   ResourceReference
	Traffic    TrafficContent
	Properties interface{} `json:"Properties,omitempty"`
}
