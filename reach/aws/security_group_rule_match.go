package aws

import (
	"fmt"
	"net"
)

type securityGroupRuleMatch struct {
	IPNetsRequired         []net.IPNet
	IP                     net.IP
	SecurityGroupReference *SecurityGroupReference
}

func matchSecurityGroupRule(
	client DomainClient,
	rule SecurityGroupRule,
	ip net.IP,
) (*securityGroupRuleMatch, error) {
	var targetIPNets []net.IPNet
	var sgRef *SecurityGroupReference

	sgRefID := rule.TargetSecurityGroupReferenceID

	if sgRefID != "" {
		var err error

		sgRef, err = client.SecurityGroupReference(sgRefID, "") // TODO: Address accountID
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
					IPNetsRequired:         nil,
					IP:                     ip,
					SecurityGroupReference: sgRef,
				}
				return match, nil
			}
		}

		return nil, nil
	}

	for _, network := range rule.TargetIPNetworks {
		if network.Contains(ip) {
			match := &securityGroupRuleMatch{
				IPNetsRequired:         targetIPNets,
				IP:                     ip,
				SecurityGroupReference: sgRef,
			}
			return match, nil
		}
	}

	return nil, nil
}

func (m securityGroupRuleMatch) Basis() securityGroupRuleMatchBasis {
	if m.SecurityGroupReference != nil {
		return securityGroupRuleMatchBasisSGRef
	}

	return securityGroupRuleMatchBasisIP
}

func (m securityGroupRuleMatch) Requirement() string {
	// TODO: Implement
	panic("implement me!")
}
