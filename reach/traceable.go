package reach

import "net"

type Traceable interface {
	Visitable(alreadyVisited bool) bool
	Destination(ips []net.IP) bool
	Segments() bool
	NextTuple(prev *IPTuple) IPTuple
	Next(t IPTuple) []InfrastructureReference
	Factors() []Factor // TODO: Need to determine what context a piece of infrastructure would need to generate this on its own (previously, this was done centrally, and involved Perspectives)
}
