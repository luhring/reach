package aws

import (
	"fmt"
	"net"
	"sort"

	"github.com/luhring/reach/reach"
)

// ResourceKindRouteTable specifies the unique name for the route table kind of resource.
const ResourceKindRouteTable reach.Kind = "RouteTable"

// A RouteTable resource representation.
type RouteTable struct {
	ID     string
	VPCID  string
	Routes []RouteTableRoute
}

// RouteTableRef returns a Reference for a RouteTable with the specified ID.
func RouteTableRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindRouteTable,
		ID:     id,
	}
}

// Resource returns the route table converted to a generalized Reach resource.
func (rt RouteTable) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindRouteTable,
		Properties: rt,
	}
}

// Ref returns a Reference for the RouteTable.
func (rt RouteTable) Ref() reach.Reference {
	return RouteTableRef(rt.ID)
}

func (rt RouteTable) routeTarget(ip net.IP) (*RouteTableRouteTarget, error) {
	routes := rt.routesBySpecificity()
	for _, route := range routes {
		if route.contains(ip) && route.State == RouteStateActive {
			return &route.Target, nil
		}
	}

	return nil, fmt.Errorf("no active RouteTableRouteTarget found for %s", ip)
}

func (rt RouteTable) routesBySpecificity() []RouteTableRoute {
	routes := rt.Routes
	sort.Sort(byRouteDestinationSpecificity(routes))
	return routes
}
