package reach

import "net"

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Ref() InfrastructureReference
	Segments() bool
	ForwardEdges(latestTuple *IPTuple, destIPs []net.IP, provider InfrastructureGetter) ([]PathEdge, error)
	Factors() []Factor
}
