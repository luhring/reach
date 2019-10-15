package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

type NetworkACLRuleAction int

const (
	NetworkACLRuleActionDeny NetworkACLRuleAction = iota
	NetworkACLRuleActionAllow
)

type NetworkACLRule struct {
	Number          int64
	TrafficContent  reach.TrafficContent
	TargetIPNetwork *net.IPNet
	Action          NetworkACLRuleAction
}
