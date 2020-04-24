package aws

import "github.com/luhring/reach/reach"

type securityGroupRulesFactorComponent struct {
	SecurityGroupID string
	RuleDirection   securityGroupRuleDirection
	RuleIndex       int
	Match           securityGroupRuleMatch
	Traffic         reach.TrafficContent
}
