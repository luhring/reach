package reach

import "net"

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	ForwardEdges(
		previousEdge *Edge,
		domains DomainProvider,
		destinationIPs []net.IP,
	) ([]Edge, error)
	Factors(
		previousEdge *Edge,
		domains DomainProvider,
	) ([]Factor, error)
}
