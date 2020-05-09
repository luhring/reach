package reach

import "net"

// Traceable is the interface that infrastructure objects can implement to be able to be traced by a tracer as points along a network path.
type Traceable interface {
	Referable
	Visitable(alreadyVisited bool) bool
	Segments() bool
	EdgesForward(resolver DomainClientResolver, previousEdge *Edge, previousRef *Reference, destinationIPs []net.IP) ([]Edge, error)
	FactorsForward(resolver DomainClientResolver, previousEdge *Edge) ([]Factor, error)
	FactorsReturn(resolver DomainClientResolver, nextEdge *Edge) ([]Factor, error)
}
