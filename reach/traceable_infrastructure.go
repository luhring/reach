package reach

import "net"

type TraceableInfrastructure interface {
	IsDestinationForIP(ips []net.IP) bool
	Next(t IPTuple) []InfrastructureReference
	UpdatedTuple(prev *IPTuple) IPTuple
	Factors() []Factor // TODO: Need to determine what context a piece of infrastructure would need to generate this on its own (previously, this was done centrally, and involved Perspectives)
	Segments() bool
}
