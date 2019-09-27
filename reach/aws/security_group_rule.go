package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

type SecurityGroupRule struct {
	TrafficContent                        reach.TrafficContent
	TargetSecurityGroupReferenceID        string       `json:"TargetSecurityGroupReferenceID,omitempty"`
	TargetSecurityGroupReferenceAccountID string       `json:"TargetSecurityGroupReferenceAccountID,omitempty"`
	TargetIPNetworks                      []*net.IPNet `json:"TargetIPNetworks,omitempty"`
}

func (rule SecurityGroupRule) MatchByIP(ip net.IP) *SecurityGroupRuleMatch {
	for _, network := range rule.TargetIPNetworks {
		if network.Contains(ip) {
			return &SecurityGroupRuleMatch{
				Basis: SecurityGroupRuleMatchBasisIP,
				Value: ip,
			}
		}
	}

	return nil
}

func (rule SecurityGroupRule) MatchBySecurityGroup(eni *ElasticNetworkInterface) *SecurityGroupRuleMatch {
	if eni != nil {
		for _, targetENISecurityGroupID := range eni.SecurityGroupIDs {
			if rule.TargetSecurityGroupReferenceID == targetENISecurityGroupID { // TODO: Handle SG Account ID
				return &SecurityGroupRuleMatch{
					Basis: SecurityGroupRuleMatchBasisSGRef,
					Value: targetENISecurityGroupID,
				}
			}
		}
	}

	return nil
}
