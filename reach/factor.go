package reach

type Factor struct {
	Kind              string
	AssociatedRole    string
	Resource          ResourceReference
	TrafficContentSet []TrafficContent
	Properties        interface{}
}
