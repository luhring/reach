package aws

import (
	"encoding/json"
	"net"
)

type networkACLRuleMatch struct {
	IP           net.IP
	MatchedIPNet net.IPNet
}

func matchNetworkACLRule(
	rule NetworkACLRule,
	ip net.IP,
) *networkACLRuleMatch {
	if rule.TargetIPNetwork.Contains(ip) {
		return &networkACLRuleMatch{
			MatchedIPNet: *rule.TargetIPNetwork,
			IP:           ip,
		}
	}

	return nil
}

func (m networkACLRuleMatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		IP           string
		MatchedIPNet string
	}{
		IP:           m.IP.String(),
		MatchedIPNet: m.MatchedIPNet.String(),
	})
}
