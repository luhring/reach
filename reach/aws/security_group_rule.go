package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// A SecurityGroupRule resource representation.
type SecurityGroupRule struct {
	TrafficContent                        reach.TrafficContent
	TargetSecurityGroupReferenceID        string      `json:"TargetSecurityGroupReferenceID,omitempty"`
	TargetSecurityGroupReferenceAccountID string      `json:"TargetSecurityGroupReferenceAccountID,omitempty"`
	TargetIPNetworks                      []net.IPNet `json:"TargetIPNetworks,omitempty"`
}
