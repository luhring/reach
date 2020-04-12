package reach

import "net"

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Destination(ips []net.IP, provider InfrastructureGetter) bool
	Segments() bool
	NextTuple(prev *IPTuple) *IPTuple
	Next(t *IPTuple, provider InfrastructureGetter) ([]InfrastructureReference, error)
	Factors() []Factor
}
