package reach

type NetworkVector struct {
	ID          string
	Source      NetworkPoint
	Destination NetworkPoint
	Factors     []Factor
}

type NetworkPoint struct {
	Domain string
	Kind   string
	ID     string
}

type Factor struct {
	Kind              string
	AssociatedRole    string
	Resource          ResourceReference
	TrafficContentSet []TrafficContent
	Properties        interface{}
}

type ResourceReference struct {
	Domain string
	Kind   string
	ID     string
}

type SecurityGroupRuleFactor struct { // (to AWS package)
	RuleIndex  int
	MatchBasis string // named type?
	MatchValue string // IP address, SG Ref ID, Prefix list name?
}
