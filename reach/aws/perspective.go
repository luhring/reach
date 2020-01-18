package aws

import "github.com/luhring/reach/reach"

type perspective struct {
	securityGroupRules                       func(sg SecurityGroup) []SecurityGroupRule
	securityGroupRuleDirection               securityGroupRuleDirection
	networkACLRulesForForwardTraffic         func(nacl NetworkACL) []NetworkACLRule
	networkACLRuleDirectionForForwardTraffic networkACLRuleDirection
	networkACLRulesForReturnTraffic          func(nacl NetworkACL) []NetworkACLRule
	networkACLRuleDirectionForReturnTraffic  networkACLRuleDirection
}

func newPerspective(p reach.Perspective) perspective {
	if p.SelfRole == reach.SubjectRoleSource {
		return newPerspectiveSourceOriented()
	}

	return newPerspectiveDestinationOriented()
}

func newPerspectiveSourceOriented() perspective {
	return perspective{
		securityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.OutboundRules
		},
		securityGroupRuleDirection: securityGroupRuleDirectionOutbound,
		networkACLRulesForForwardTraffic: func(nacl NetworkACL) []NetworkACLRule {
			return nacl.OutboundRules
		},
		networkACLRuleDirectionForForwardTraffic: networkACLRuleDirectionOutbound,
		networkACLRulesForReturnTraffic: func(nacl NetworkACL) []NetworkACLRule {
			return nacl.InboundRules
		},
		networkACLRuleDirectionForReturnTraffic: networkACLRuleDirectionInbound,
	}
}

func newPerspectiveDestinationOriented() perspective {
	return perspective{
		securityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.InboundRules
		},
		securityGroupRuleDirection: securityGroupRuleDirectionInbound,
		networkACLRulesForForwardTraffic: func(nacl NetworkACL) []NetworkACLRule {
			return nacl.InboundRules
		},
		networkACLRuleDirectionForForwardTraffic: networkACLRuleDirectionInbound,
		networkACLRulesForReturnTraffic: func(nacl NetworkACL) []NetworkACLRule {
			return nacl.OutboundRules
		},
		networkACLRuleDirectionForReturnTraffic: networkACLRuleDirectionOutbound,
	}
}
