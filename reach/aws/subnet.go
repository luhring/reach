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

// SubnetRef returns a Reference for a Subnet with the specified ID.
func SubnetRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSubnet,
		ID:     id,
	}
}

// Resource returns the subnet converted to a generalized Reach resource.
func (s Subnet) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindSubnet,
		Properties: s,
	}
}

// Ref returns a Reference for the Subnet.
func (s Subnet) Ref() reach.Reference {
	return SubnetRef(s.ID)
}

// equal returns a boolean to indicate if two Subnet instances represent the same Subnet.
func (s Subnet) equal(other Subnet) bool {
	return s.ID == other.ID
}

// contains returns a boolean to indicate if the specified IP address is contained within any of the Subnet's CIDR blocks.
func (s Subnet) contains(ip net.IP) bool {
	if s.IPv4CIDR.Contains(ip) {
		return true
	}

	if s.IPv6CIDR != nil && s.IPv6CIDR.Contains(ip) {
		return true
	}

	return false
}
