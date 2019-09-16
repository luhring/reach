package aws

import "net"

type RouteTableRoute struct {
	Destination *net.IPNet  `json:"destination"`
	Target      interface{} `json:"target"` // TODO: Figure this out -- this is not the normal Reach 'target'
	States      string      `json:"states"`
	Propagated  bool        `json:"propagated"`
}
