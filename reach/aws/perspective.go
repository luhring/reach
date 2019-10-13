package aws

type Perspective struct {
	getSecurityGroupRules func(sg SecurityGroup) []SecurityGroupRule
	ruleDirection         SecurityGroupRuleDirection
}

func NewPerspectiveSourceOriented() Perspective {
	return Perspective{
		getSecurityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.OutboundRules
		},
		ruleDirection: SecurityGroupRuleDirectionOutbound,
	}
}

func NewPerspectiveDestinationOriented() Perspective {
	return Perspective{
		getSecurityGroupRules: func(sg SecurityGroup) []SecurityGroupRule {
			return sg.InboundRules
		},
		ruleDirection: SecurityGroupRuleDirectionInbound,
	}
}
