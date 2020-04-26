package reach

import "net"

type Traceable interface {
	Ref() UniversalReference
	Visitable(alreadyVisited bool) bool
	Segments() bool
	EdgesForward(resolver DomainClientResolver, previousEdge *Edge, previousRef *UniversalReference, destinationIPs []net.IP) ([]Edge, error)
	FactorsForward(resolver DomainClientResolver, previousEdge *Edge) ([]Factor, error)
	FactorsReturn(resolver DomainClientResolver, nextEdge *Edge) ([]Factor, error)
}
