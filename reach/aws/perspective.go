package aws

type perspective struct {
	getSecurityGroupRules func(sg SecurityGroup) []SecurityGroupRule
	ruleDirection         securityGroupRuleDirection
}

func newPerspectiveSourceOriented() perspective {
	return perspective{
		getSecurityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.OutboundRules
		},
		ruleDirection: securityGroupRuleDirectionOutbound,
	}
}

func newPerspectiveDestinationOriented() perspective {
	return perspective{
		getSecurityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.InboundRules
		},
		ruleDirection: securityGroupRuleDirectionInbound,
	}
}
