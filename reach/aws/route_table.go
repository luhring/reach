package aws

import "github.com/luhring/reach/reach"

// ResourceKindRouteTable specifies the unique name for the route table kind of resource.
const ResourceKindRouteTable reach.Kind = "RouteTable"

// A RouteTable resource representation.
type RouteTable struct {
	ID     string
	VPCID  string
	Routes []RouteTableRoute
}

// Resource returns the route table converted to a generalized Reach resource.
func (rt RouteTable) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindRouteTable,
		Properties: rt,
	}
}
