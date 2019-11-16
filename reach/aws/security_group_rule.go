package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// A SecurityGroupRule resource representation.
type SecurityGroupRule struct {
	TrafficContent                        reach.TrafficContent
	TargetSecurityGroupReferenceID        string       `json:"TargetSecurityGroupReferenceID,omitempty"`
	TargetSecurityGroupReferenceAccountID string       `json:"TargetSecurityGroupReferenceAccountID,omitempty"`
	TargetIPNetworks                      []*net.IPNet `json:"TargetIPNetworks,omitempty"`
}

func (rule SecurityGroupRule) matchByIP(ip net.IP) *securityGroupRuleMatch {
	for _, network := range rule.TargetIPNetworks {
		if network.Contains(ip) {
			return &securityGroupRuleMatch{
				Basis:       securityGroupRuleMatchBasisIP,
				Requirement: network,
				Value:       ip,
			}
		}
	}

	return nil
}

func (rule SecurityGroupRule) matchBySecurityGroup(eni *ElasticNetworkInterface) *securityGroupRuleMatch {
	if eni != nil {
		for _, targetENISecurityGroupID := range eni.SecurityGroupIDs {
			if rule.TargetSecurityGroupReferenceID == targetENISecurityGroupID { // TODO: Handle SG Account ID
				return &securityGroupRuleMatch{
					Basis:       securityGroupRuleMatchBasisSGRef,
					Requirement: rule.TargetSecurityGroupReferenceID,
					Value:       targetENISecurityGroupID,
				}
			}
		}
	}

	return nil
}
