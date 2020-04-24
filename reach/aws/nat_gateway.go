package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindNATGateway specifies the unique name for the NAT gateway kind of resource.
const ResourceKindNATGateway reach.Kind = "NATGateway"

// A NATGateway resource representation.
type NATGateway struct {
	ID        string
	SubnetID  string
	VPCID     string
	PrivateIP net.IP
	PublicIP  net.IP
}

// Resource returns the NAT gateway converted to a generalized Reach resource.
func (ngw NATGateway) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindNATGateway,
		Properties: ngw,
	}
}
