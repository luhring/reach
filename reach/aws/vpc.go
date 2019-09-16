package aws

import "net"

type VPC struct {
	ID        string      `json:"id"`
	IPv4CIDRs []net.IPNet `json:"ipv4CIDRs"`
	IPv6CIDRs []net.IPNet `json:"ipv6CIDRs"`
}
