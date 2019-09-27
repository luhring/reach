package reach

type Factor struct {
	Kind       string
	Resource   ResourceReference
	Traffic    TrafficContent
	Properties interface{} `json:"Properties,omitempty"`
}
