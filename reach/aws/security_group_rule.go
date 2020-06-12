package aws

import (
	"net"

	"github.com/luhring/reach/reach/traffic"
)

// A SecurityGroupRule resource representation.
type SecurityGroupRule struct {
	TrafficContent                        traffic.Content
	TargetSecurityGroupReferenceID        string      `json:"TargetSecurityGroupReferenceID,omitempty"`
	TargetSecurityGroupReferenceAccountID string      `json:"TargetSecurityGroupReferenceAccountID,omitempty"`
	TargetIPNetworks                      []net.IPNet `json:"TargetIPNetworks,omitempty"`
}
