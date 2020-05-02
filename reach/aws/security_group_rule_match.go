package aws

import (
	"fmt"
	"net"
)

type securityGroupRuleMatch struct {
	IP           net.IP
	MatchedIPNet *net.IPNet
	MatchedSGRef *SecurityGroupReference
}

func matchSecurityGroupRule(
	client DomainClient,
	rule SecurityGroupRule,
	ip net.IP,
) (*securityGroupRuleMatch, error) {
	if sgRefID := rule.TargetSecurityGroupReferenceID; sgRefID != "" {
		var err error

		sgRef, err := client.SecurityGroupReference(sgRefID, "") // TODO: Address accountID
		if err != nil {
			return nil, fmt.Errorf("unable to determine rule match: %v", err)
		}

		enis, err := client.ResolveSecurityGroupReference(sgRefID)
		if err != nil {
			return nil, fmt.Errorf("unable to determine rule match: %v", err)
		}

		for _, eni := range enis {
			if eni.owns(ip) {
				match := &securityGroupRuleMatch{
					IP:           ip,
					MatchedIPNet: nil,
					MatchedSGRef: sgRef,
				}
				return match, nil
			}
		}

		return nil, nil
	}

	for _, network := range rule.TargetIPNetworks {
		if network.Contains(ip) {
			match := &securityGroupRuleMatch{
				IP:           ip,
				MatchedIPNet: &network,
				MatchedSGRef: nil,
			}
			return match, nil
		}
	}

	return nil, nil
}

func (m securityGroupRuleMatch) Basis() securityGroupRuleMatchBasis {
	if m.MatchedSGRef != nil {
		return securityGroupRuleMatchBasisSGRef
	}

	return securityGroupRuleMatchBasisIP
}

func (m securityGroupRuleMatch) Requirement() string {
	// TODO: Implement
	panic("implement me!")
}
