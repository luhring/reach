package aws

import "github.com/luhring/reach/reach"

type securityGroupRulesFactorComponent struct {
	SecurityGroup reach.ResourceReference
	RuleDirection securityGroupRuleDirection
	RuleIndex     int
	Match         securityGroupRuleMatch
	Traffic       reach.TrafficContent
}
