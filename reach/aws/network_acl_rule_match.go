package aws

import "net"

type networkACLRuleMatch struct {
	IPNetRequired net.IPNet
	IP            net.IP
}

func matchNetworkACLRule(
	rule NetworkACLRule,
	ip net.IP,
) *networkACLRuleMatch {
	if rule.TargetIPNetwork.Contains(ip) {
		return &networkACLRuleMatch{
			IPNetRequired: *rule.TargetIPNetwork,
			IP:            ip,
		}
	}

	return nil
}
