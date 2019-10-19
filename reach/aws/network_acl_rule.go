package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

// A NetworkACLRuleAction is the action specified by a network ACL rule -- either allow or deny.
type NetworkACLRuleAction int

// The allowed actions for a network ACL rule.
const (
	NetworkACLRuleActionDeny NetworkACLRuleAction = iota
	NetworkACLRuleActionAllow
)

// An NetworkACLRule resource representation.
type NetworkACLRule struct {
	Number          int64
	TrafficContent  reach.TrafficContent
	TargetIPNetwork *net.IPNet
	Action          NetworkACLRuleAction
}
