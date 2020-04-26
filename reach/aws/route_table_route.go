package aws

import (
	"net"
)

// A RouteTableRoute resource representation.
type RouteTableRoute struct {
	Destination *net.IPNet
	State       RouteTableRouteState
	Target      RouteTableRouteTarget
}

func (route RouteTableRoute) maskZeros() int {
	ones, bits := route.Destination.Mask.Size()
	return bits - ones
}

func (route RouteTableRoute) contains(ip net.IP) bool {
	return route.Destination.Contains(ip)
}

type byRouteDestinationSpecificity []RouteTableRoute

func (s byRouteDestinationSpecificity) Len() int {
	return len(s)
}

func (s byRouteDestinationSpecificity) Less(i, j int) bool {
	return s[i].maskZeros() < s[j].maskZeros()
}

func (s byRouteDestinationSpecificity) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
