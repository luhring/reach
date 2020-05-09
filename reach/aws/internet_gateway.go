package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindInternetGateway specifies the unique name for the InternetGateway kind of resource.
const ResourceKindInternetGateway reach.Kind = "InternetGateway"

// An InternetGateway resource representation.
type InternetGateway struct {
	ID    string
	VPCID string
}

// InternetGatewayRef returns a Reference for an InternetGateway with the specified ID.
func InternetGatewayRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindInternetGateway,
		ID:     id,
	}
}

// Resource returns the Internet gateway converted to a generalized Reach resource.
func (igw InternetGateway) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindInternetGateway,
		Properties: igw,
	}
}

// ———— Implementing Traceable ————

// Ref returns a Reference for the InternetGateway.
func (igw InternetGateway) Ref() reach.Reference {
	return InternetGatewayRef(igw.ID)
}

// Visitable returns a boolean to indicate whether a tracer is allowed to add this resource to the path it's currently constructing.
func (igw InternetGateway) Visitable(alreadyVisited bool) bool {
	panic("implement me")
}

// Segments returns a boolean to indicate whether a tracer should create a new path segment at this point in the path.
func (igw InternetGateway) Segments() bool {
	panic("implement me")
}

// EdgesForward returns the set of all possible edges forward given this point in a path that a tracer is constructing. EdgesForward returns an empty slice of edges if there are no further points for the specified network traffic to travel as it attempts to reach its intended network destination.
func (igw InternetGateway) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, previousRef *reach.Reference, destinationIPs []net.IP) ([]reach.Edge, error) {
	panic("implement me")
}

// FactorsForward returns a set of factors that impact the traffic traveling through this point in the direction of source to destination.
func (igw InternetGateway) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

// FactorsReturn returns a set of factors that impact the traffic traveling through this point in the direction of destination to source.
func (igw InternetGateway) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}
