package aws

import "github.com/luhring/reach/reach"

type networkACLRulesFactorComponent struct {
	NetworkACL    reach.ResourceReference
	RuleDirection networkACLRuleDirection
	RuleNumber    int64
	Match         networkACLRuleMatch
	Traffic       reach.TrafficContent
}
