package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindSubnet specifies the unique name for the subnet kind of resource.
const ResourceKindSubnet reach.Kind = "Subnet"

// A Subnet resource representation.
type Subnet struct {
	ID           string
	NetworkACLID string
	RouteTableID string
	VPCID        string
	IPv4CIDR     net.IPNet
	IPv6CIDR     *net.IPNet
}

// Resource returns the subnet converted to a generalized Reach resource.
func (s Subnet) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSubnet,
		Properties: s,
	}
}

func (s Subnet) equal(other Subnet) bool {
	return s.ID == other.ID
}
