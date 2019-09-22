package reach

type Factor struct {
	Kind              string
	Resource          ResourceReference
	TrafficContentSet []TrafficContent
	Properties        interface{}
}
