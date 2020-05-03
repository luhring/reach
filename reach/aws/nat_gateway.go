package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindNATGateway specifies the unique name for the NAT gateway kind of resource.
const ResourceKindNATGateway reach.Kind = "NATGateway"

// A NATGateway resource representation.
type NATGateway struct {
	ID        string
	SubnetID  string
	VPCID     string
	PrivateIP net.IP
	PublicIP  net.IP
}

// Resource returns the NAT gateway converted to a generalized Reach resource.
func (ngw NATGateway) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNATGateway,
		Properties: ngw,
	}
}

func (ngw NATGateway) Ref() reach.UniversalReference {
	return NATGatewayRef(ngw.ID)
}

func (ngw NATGateway) Visitable(alreadyVisited bool) bool {
	panic("implement me")
}

func (ngw NATGateway) Segments() bool {
	panic("implement me")
}

func (ngw NATGateway) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, previousRef *reach.UniversalReference, destinationIPs []net.IP) ([]reach.Edge, error) {
	panic("implement me")
}

func (ngw NATGateway) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

func (ngw NATGateway) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

func NATGatewayRef(id string) reach.UniversalReference {
	return reach.UniversalReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindNATGateway,
		ID:     id,
	}
}
