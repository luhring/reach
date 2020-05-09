package aws

import "github.com/luhring/reach/reach"

type networkACLRulesFactorComponent struct {
	NetworkACLID  string
	RuleDirection networkACLRuleDirection
	RuleIndex     int64
	Match         networkACLRuleMatch
	Traffic       reach.TrafficContent
}