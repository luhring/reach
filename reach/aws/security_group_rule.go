package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

type SecurityGroupRule struct {
	ProtocolContent                       reach.ProtocolContent `json:"protocolContent"`
	TargetSecurityGroupReferenceID        string                `json:"targetSecurityGroupReferenceID"`
	TargetSecurityGroupReferenceAccountID string                `json:"targetSecurityGroupReferenceAccountID"`
	TargetIPNetworks                      []*net.IPNet          `json:"targetIPNetworks"`
}
