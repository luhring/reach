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
	resources DomainClient,
	rule SecurityGroupRule,
	ip net.IP,
) (*securityGroupRuleMatch, error) {
	var targetIPNets []net.IPNet
	var sgRef *SecurityGroupReference

	sgRefID := rule.TargetSecurityGroupReferenceID
	if sgRefID != "" {
		var err error
		targetIPNets, err = resources.ResolveSecurityGroupReference(sgRefID)
		if err != nil {
			return nil, fmt.Errorf("unable to determine rule match: %v", err)
		}
		sgRef, err = resources.SecurityGroupReference(sgRefID, "") // TODO: Address accountID
		if err != nil {
			return nil, fmt.Errorf("unable to determine rule match: %v", err)
		}
	} else {
		targetIPNets = rule.TargetIPNetworks
	}

	for _, network := range targetIPNets {
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
