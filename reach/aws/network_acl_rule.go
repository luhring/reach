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

// Allows returns a boolean indicating if the rule is allowing traffic.
func (r NetworkACLRule) Allows() bool {
	return r.Action == NetworkACLRuleActionAllow
}

// Denies returns a boolean indicating if the rule is denying traffic.
func (r NetworkACLRule) Denies() bool {
	return r.Action == NetworkACLRuleActionDeny
}

func (r NetworkACLRule) matchByIP(ip net.IP) *networkACLRuleMatch {
	if r.TargetIPNetwork.Contains(ip) {
		return &networkACLRuleMatch{
			Requirement: *r.TargetIPNetwork,
			Value:       ip,
		}
	}

	return nil
}
