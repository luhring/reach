package reach

import "net"

type Traceable interface {
	Referable
	Visitable(alreadyVisited bool) bool
	Segments() bool
	EdgesForward(resolver DomainClientResolver, previousEdge *Edge, previousRef *Reference, destinationIPs []net.IP) ([]Edge, error)
	FactorsForward(resolver DomainClientResolver, previousEdge *Edge) ([]Factor, error)
	FactorsReturn(resolver DomainClientResolver, nextEdge *Edge) ([]Factor, error)
}
