package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindVPC specifies the unique name for the VPC kind of resource.
const ResourceKindVPC = "VPC"

// An VPC resource representation.
type VPC struct {
	ID        string
	IPv4CIDRs []net.IPNet `json:"IPv4CIDRs,omitempty"`
	IPv6CIDRs []net.IPNet `json:"IPv6CIDRs,omitempty"`
}

// ToResource returns the VPC converted to a generalized Reach resource.
func (vpc VPC) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPC,
		Properties: vpc,
	}
}
