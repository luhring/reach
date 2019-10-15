package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindVPC = "VPC"

type VPC struct {
	ID        string
	IPv4CIDRs []net.IPNet `json:"IPv4CIDRs,omitempty"`
	IPv6CIDRs []net.IPNet `json:"IPv6CIDRs,omitempty"`
}

func (vpc VPC) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPC,
		Properties: vpc,
	}
}
