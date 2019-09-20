package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

const ResourceKindVPC = "VPC"

type VPC struct {
	ID        string      `json:"id"`
	IPv4CIDRs []net.IPNet `json:"ipv4CIDRs"`
	IPv6CIDRs []net.IPNet `json:"ipv6CIDRs"`
}

func (vpc VPC) ToResource() reach.Resource {
	return reach.Resource{
		Kind:       ResourceKindVPC,
		Properties: vpc,
	}
}
