package aws

import (
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
)

type VPCRouter struct {
	VPC VPC
}

func NewVPCRouter(client DomainClient, id string) (*VPCRouter, error) {
	vpc, err := client.VPC(id)
	if err != nil {
		return nil, fmt.Errorf("unable to get VPC: %v", err)
	}

	return &VPCRouter{VPC: *vpc}, nil
}

// ———— Implementing Traceable ————

func (r VPCRouter) Ref() reach.UniversalReference {
	return reach.UniversalReference{
		Implicit: true,
		R:        r.VPC.ResourceReference(),
	}
}

func (r VPCRouter) Visitable(alreadyVisited bool) bool {
	panic("implement me")
}

func (r VPCRouter) Segments() bool {
	panic("implement me")
}

func (r VPCRouter) EdgesForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge, destinationIPs []net.IP) ([]reach.Edge, error) {
	panic("implement me")
}

func (r VPCRouter) FactorsForward(resolver reach.DomainClientResolver, previousEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}

func (r VPCRouter) FactorsReturn(resolver reach.DomainClientResolver, nextEdge *reach.Edge) ([]reach.Factor, error) {
	panic("implement me")
}
