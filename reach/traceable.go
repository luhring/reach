package reach

import "net"

type Traceable interface {
	AllowsVisit(previouslyVisited bool) bool
	IsDestinationForIP(ips []net.IP) bool
	UpdatedTuple(prev *IPTuple) IPTuple
	Factors() []Factor // TODO: Need to determine what context a piece of infrastructure would need to generate this on its own (previously, this was done centrally, and involved Perspectives)
	Segments() bool
	Next(t IPTuple) []InfrastructureReference
}
