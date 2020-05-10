package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindNATGateway specifies the unique name for the NAT gateway kind of resource.
const ResourceKindNATGateway reach.Kind = "NATGateway"

var _ reach.Traceable = (*NATGateway)(nil)

// A NATGateway resource representation.
type NATGateway struct {
	ID        string
	SubnetID  string
	VPCID     string
	PrivateIP net.IP
	PublicIP  net.IP
}

// NATGatewayRef returns a Reference for a NATGateway with the specified ID.
func NATGatewayRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindNATGateway,
		ID:     id,
	}
}

// Resource returns the NATGateway converted to a generalized Reach resource.
func (ngw NATGateway) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNATGateway,
		Properties: ngw,
	}
}

// ———— Implementing Traceable ————

// Ref returns a Reference for the NATGateway.
func (ngw NATGateway) Ref() reach.Reference {
	return NATGatewayRef(ngw.ID)
}

// Visitable returns a boolean to indicate whether a tracer is allowed to add this resource to the path it's currently constructing.
func (ngw NATGateway) Visitable(alreadyVisited bool) bool {
	panic("implement me")
}

// Segments returns a boolean to indicate whether a tracer should create a new path segment at this point in the path.
func (ngw NATGateway) Segments() bool {
	panic("implement me")
}

// EdgesForward returns the set of all possible edges forward given this point in a path that a tracer is constructing. EdgesForward returns an empty slice of edges if there are no further points for the specified network traffic to travel as it attempts to reach its intended network destination.
func (ngw NATGateway) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, previousRef *reach.Reference, destinationIPs []net.IP) ([]reach.Edge, error) {
	panic("implement me")
}

// FactorsForward returns a set of factors that impact the traffic traveling through this point in the direction of source to destination.
func (ngw NATGateway) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

// FactorsReturn returns a set of factors that impact the traffic traveling through this point in the direction of destination to source.
func (ngw NATGateway) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}
