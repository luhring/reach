package aws

import "net"

type RouteTableRoute struct {
	Destination *net.IPNet
	Target      interface{} // TODO: Figure this out -- this is not the normal Reach 'target'
	States      string
	Propagated  bool
}
