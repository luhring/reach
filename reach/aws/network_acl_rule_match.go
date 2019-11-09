package aws

import "net"

type networkACLRuleMatch struct {
	Requirement net.IPNet
	Value       net.IP
}
