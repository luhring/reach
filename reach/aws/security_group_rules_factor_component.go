package aws

import (
	"github.com/luhring/reach/reach/traffic"
)

type securityGroupRulesFactorComponent struct {
	SecurityGroupID string
	RuleDirection   securityGroupRuleDirection
	RuleIndex       int
	Match           securityGroupRuleMatch
	Traffic         traffic.Content
}
