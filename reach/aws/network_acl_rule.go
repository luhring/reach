package aws

import (
	"encoding/json"
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

// String returns the string representation of the NetworkACLRuleAction.
func (action NetworkACLRuleAction) String() string {
	switch action {
	case NetworkACLRuleActionDeny:
		return "deny"
	case NetworkACLRuleActionAllow:
		return "allow"
	default:
		return "[unknown action]"
	}
}

// MarshalJSON returns the JSON representation of the NetworkACLRuleAction.
func (action NetworkACLRuleAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(action.String())
}

// An NetworkACLRule resource representation.
type NetworkACLRule struct {
	Number          int64
	TrafficContent  reach.TrafficContent
	TargetIPNetwork *net.IPNet
	Action          NetworkACLRuleAction
}

func (r NetworkACLRule) Allows() bool {
	return r.Action == NetworkACLRuleActionAllow
}

func (r NetworkACLRule) Denies() bool {
	return r.Action == NetworkACLRuleActionDeny
}

func (rule NetworkACLRule) matchByIP(ip net.IP) *networkACLRuleMatch {
	if rule.TargetIPNetwork.Contains(ip) {
		return &networkACLRuleMatch{
			Requirement: *rule.TargetIPNetwork,
			Value:       ip,
		}
	}

	return nil
}
