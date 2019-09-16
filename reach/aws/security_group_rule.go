package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

type SecurityGroupRule struct {
	TrafficContent                        reach.TrafficContent `json:"trafficContent"`
	TargetSecurityGroupReferenceID        string               `json:"targetSecurityGroupReferenceID"`
	TargetSecurityGroupReferenceAccountID string               `json:"targetSecurityGroupReferenceAccountID"`
	TargetIPNetworks                      []*net.IPNet         `json:"targetIPNetworks"`
}
