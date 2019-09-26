package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

type SecurityGroupRule struct {
	TrafficContent                        reach.TrafficContent `json:"trafficContent"`
	TargetSecurityGroupReferenceID        string               `json:"targetSecurityGroupReferenceID,omitempty"`
	TargetSecurityGroupReferenceAccountID string               `json:"targetSecurityGroupReferenceAccountID,omitempty"`
	TargetIPNetworks                      []*net.IPNet         `json:"targetIPNetworks,omitempty"`
}
