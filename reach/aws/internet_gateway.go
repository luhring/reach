package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindInternetGateway specifies the unique name for the Internet gateway kind of resource.
const ResourceKindInternetGateway reach.Kind = "InternetGateway"

// An InternetGateway resource representation.
type InternetGateway struct {
	ID    string
	VPCID string
}

// Resource returns the Internet gateway converted to a generalized Reach resource.
func (igw InternetGateway) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindInternetGateway,
		Properties: igw,
	}
}

func (igw InternetGateway) Ref() reach.UniversalReference {
	return InternetGatewayRef(igw.ID)
}

func (igw InternetGateway) Visitable(alreadyVisited bool) bool {
	panic("implement me")
}

func (igw InternetGateway) Segments() bool {
	panic("implement me")
}

func (igw InternetGateway) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, previousRef *reach.UniversalReference, destinationIPs []net.IP) ([]reach.Edge, error) {
	panic("implement me")
}

func (igw InternetGateway) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

func (igw InternetGateway) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

func InternetGatewayRef(id string) reach.UniversalReference {
	return reach.UniversalReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindInternetGateway,
		ID:     id,
	}
}
