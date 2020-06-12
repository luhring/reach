package aws

import (
	"github.com/luhring/reach/reach/traffic"
)

type networkACLRulesFactorComponent struct {
	NetworkACLID  string
	RuleDirection NetworkACLRuleDirection
	RuleIndex     int64
	Match         networkACLRuleMatch
	Traffic       traffic.Content
}
