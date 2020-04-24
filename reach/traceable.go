package reach

import "net"

type Traceable interface {
	Ref() UniversalReference
	Visitable(alreadyVisited bool) bool
	Segments() bool
	EdgesForward(domains DomainProvider, previousEdge *Edge, destinationIPs []net.IP) ([]Edge, error)
	FactorsForward(domains DomainProvider, previousEdge *Edge) ([]Factor, error)
	FactorsReturn(domains DomainProvider, nextEdge *Edge) ([]Factor, error)
}
