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

// VPCRef returns a Reference for a VPC with the specified ID.
func VPCRef(id string) reach.Reference {
	return reach.Reference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindVPC,
		ID:     id,
	}
}

// Resource returns the VPC converted to a generalized Reach resource.
func (vpc VPC) Resource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPC,
		Properties: vpc,
	}
}

// Ref returns a Reference for the VPC.
func (vpc VPC) Ref() reach.Reference {
	return VPCRef(vpc.ID)
}

// contains returns a boolean to indicate whether the specified IP address is contained within any of the VPC's associated IP CIDR blocks.
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

// subnetThatContains returns the Subnet that exists within the VPC, where the Subnet's CIDR block contains the specified IP address. If no such Subnet exists, subnetThatContains returns false for its second return parameter. If an error is encountered while determine which Subnet contains the IP address, the error is returned.
func (vpc VPC) subnetThatContains(client DomainClient, ip net.IP) (*Subnet, bool, error) {
	if vpc.contains(ip) == false {
		return nil, false, nil
	}

	subnets, err := client.SubnetsByVPC(vpc.ID)
	if err != nil {
		return nil, false, err
	}

	for _, s := range subnets {
		if s.contains(ip) {
			return &s, true, nil
		}
	}

	return nil, false, nil
}
