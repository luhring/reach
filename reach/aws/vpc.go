package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// ResourceKindVPC specifies the unique name for the VPC kind of resource.
const ResourceKindVPC reach.Kind = "VPC"

// An VPC resource representation.
type VPC struct {
	ID        string
	IPv4CIDRs []net.IPNet `json:"IPv4CIDRs,omitempty"`
	IPv6CIDRs []net.IPNet `json:"IPv6CIDRs,omitempty"`
}

// Resource returns the VPC converted to a generalized Reach resource.
func (vpc VPC) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPC,
		Properties: vpc,
	}
}

func (vpc VPC) ResourceReference() reach.ResourceReference {
	return reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     vpc.ID,
	}
}

func (vpc VPC) contains(ip net.IP) bool {
	for _, network := range vpc.IPv4CIDRs {
		if network.Contains(ip) {
			return true
		}
	}

	for _, network := range vpc.IPv6CIDRs {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func (vpc VPC) subnetThatContains(client DomainClient, ip net.IP) (*Subnet, error) {
	panic("implement me!")
}
