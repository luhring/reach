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
	Number          int64                `json:"number"`
	TrafficContent  reach.TrafficContent `json:"trafficContent"`
	TargetIPNetwork *net.IPNet           `json:"targetIPNetwork"`
	Action          NetworkACLRuleAction `json:"action"`
}
