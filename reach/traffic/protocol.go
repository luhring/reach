package traffic

import (
	"encoding/json"
	"fmt"
)

// A Protocol represents an analyzable IP protocol, whose integer value corresponds to the officially assigned number for the IP protocol (as defined here: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
type Protocol int

// Protocol numbers for the most common IP protocols.
const (
	ProtocolAll    Protocol = -1
	ProtocolICMPv4 Protocol = 1
	ProtocolTCP    Protocol = 6
	ProtocolUDP    Protocol = 17
	ProtocolICMPv6 Protocol = 58
)

// Names of the most common IP protocols.
const (
	ProtocolNameAll    = "all"
	ProtocolNameICMPv4 = "ICMPv4"
	ProtocolNameTCP    = "TCP"
	ProtocolNameUDP    = "UDP"
	ProtocolNameICMPv6 = "ICMPv6"
)

// UsesPorts returns a boolean indicating whether or not the described protocol has a "ports" concept that can be further drilled into when analyzing network rules. UsesPorts returns true if the underlying protocol is either TCP or UDP.
func (p Protocol) UsesPorts() bool {
	return p == ProtocolTCP || p == ProtocolUDP
}

// UsesICMPTypeCodes returns a boolean indicating whether or not the underlying protocol is ICMP (v4) or ICMPv6.
func (p Protocol) UsesICMPTypeCodes() bool {
	return p == ProtocolICMPv4 || p == ProtocolICMPv6
}

// IsCustomProtocol returns a boolean indicating whether or not the underlying protocol is "custom", meaning that it's not TCP, UDP, ICMPv4, or ICMPv6. The significance of this distinction is that Reach can analyze custom protocols only on an "all-or-nothing" basis, in contrast to protocols like TCP, where Reach can further assess traffic flow on a more granular basis, like ports.
func (p Protocol) IsCustomProtocol() bool {
	return p != ProtocolICMPv4 && p != ProtocolTCP && p != ProtocolUDP && p != ProtocolICMPv6
}

// String returns the common name of the IP protocol.
func (p Protocol) String() string {
	return ProtocolName(p)
}

// MarshalJSON returns the JSON representation of the Protocol.
func (p Protocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

// DisplayOrder returns an opinionated display order for the given protocol.
func (p Protocol) DisplayOrder() int {
	switch p {
	case ProtocolTCP:
		return -100
	case ProtocolUDP:
		return -90
	case ProtocolICMPv4:
		return -80
	case ProtocolICMPv6:
		return -70
	default:
		return int(p)
	}
}

// ProtocolName returns the name of an IP protocol given the protocol's assigned number.
func ProtocolName(protocol Protocol) string {
	switch protocol {
	case ProtocolICMPv4:
		return ProtocolNameICMPv4
	case ProtocolTCP:
		return ProtocolNameTCP
	case ProtocolUDP:
		return ProtocolNameUDP
	case ProtocolICMPv6:
		return ProtocolNameICMPv6
	default:
		return customProtocolName(protocol)
	}
}

func customProtocolName(protocol Protocol) string {
	name, exists := ipProtocols[protocol]
	if exists {
		return name
	}

	return fmt.Sprintf("IP protocol %d", protocol)
}
