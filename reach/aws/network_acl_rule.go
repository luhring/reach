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
	Number          int64                 `json:"number"`
	ProtocolContent reach.ProtocolContent `json:"protocolContent"`
	TargetIPNetwork *net.IPNet            `json:"targetIPNetwork"`
	Action          NetworkACLRuleAction  `json:"action"`
}
