package reach

import (
	"fmt"
)

type Protocol int

const (
	ProtocolAll        Protocol = -1
	ProtocolICMPv4     Protocol = 1
	ProtocolTCP        Protocol = 6
	ProtocolUDP        Protocol = 17
	ProtocolICMPv6     Protocol = 58
	ProtocolNameAll             = "all"
	ProtocolNameICMP            = "ICMP"
	ProtocolNameTCP             = "TCP"
	ProtocolNameUDP             = "UDP"
	ProtocolNameICMPv6          = "ICMPv6"
)

func (p Protocol) UsesPorts() bool {
	return p == ProtocolTCP || p == ProtocolUDP
}

func (p Protocol) UsesICMPTypeCodes() bool {
	return p == ProtocolICMPv4 || p == ProtocolICMPv6
}

func (p Protocol) IsCustomProtocol() bool {
	return p != ProtocolICMPv4 && p != ProtocolTCP && p != ProtocolUDP && p != ProtocolICMPv6
}

func ProtocolName(protocol Protocol) string {
	switch protocol {
	case ProtocolICMPv4:
		return ProtocolNameICMP
	case ProtocolTCP:
		return ProtocolNameTCP
	case ProtocolUDP:
		return ProtocolNameUDP
	case ProtocolICMPv6:
		return ProtocolNameICMPv6
	default:
		return fmt.Sprintf("[IP Protocol %d]", protocol)
	}
}
