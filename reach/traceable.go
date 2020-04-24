package reach

import "net"

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	EdgesForward(
		previousEdge *Edge,
		domains DomainProvider,
		destinationIPs []net.IP,
	) ([]Edge, error)
	FactorsForward(
		previousEdge *Edge,
		domains DomainProvider,
	) ([]Factor, error)
	FactorsReturn(
		nextEdge *Edge,
		domains DomainProvider,
	) ([]Factor, error)
}
