package reach

import (
	"fmt"
	"strings"

	"github.com/luhring/reach/reach/set"
)

type ProtocolContent struct {
	Protocol Protocol
	Ports    *set.PortSet `json:"Ports,omitempty"`
	ICMP     *set.ICMPSet `json:"ICMP,omitempty"`
}

func NewProtocolContent(protocol Protocol, ports *set.PortSet, icmp *set.ICMPSet) ProtocolContent {
	return ProtocolContent{
		protocol,
		ports,
		icmp,
	}
}

func NewProtocolContentWithPorts(protocol Protocol, ports *set.PortSet) ProtocolContent {
	return NewProtocolContent(protocol, ports, nil)
}

func NewProtocolContentWithICMP(protocol Protocol, icmp *set.ICMPSet) ProtocolContent {
	return NewProtocolContent(protocol, nil, icmp)
}

func NewProtocolContentForCustomProtocol(protocol Protocol) ProtocolContent {
	return NewProtocolContent(protocol, nil, nil)
}

func NewProtocolContentForAllTraffic() ProtocolContent {
	return NewProtocolContent(ProtocolAll, nil, nil)
}

func NewProtocolContentForNoContent() ProtocolContent {
	return NewProtocolContent(ProtocolNone, nil, nil)
}

func (pc ProtocolContent) allProtocols() bool {
	return pc.Protocol == ProtocolAll
}

func (pc ProtocolContent) noContent() bool {
	return pc.Protocol == ProtocolNone
}

func (pc ProtocolContent) isTCPOrUDP() bool {
	return pc.Protocol == ProtocolTCP || pc.Protocol == ProtocolUDP
}

func (pc ProtocolContent) isICMPv4OrICMPv6() bool {
	return pc.Protocol == ProtocolICMPv4 || pc.Protocol == ProtocolICMPv6
}

func (pc ProtocolContent) intersect(other ProtocolContent) (ProtocolContent, error) {
	if pc.noContent() || other.noContent() {
		return NewProtocolContentForNoContent(), nil
	}

	if pc.allProtocols() {
		return other, nil
	}

	if other.allProtocols() {
		return pc, nil
	}

	if pc.sameProtocolAs(other) == false {
		return ProtocolContent{}, fmt.Errorf(
			"cannot intersect with different protocols (IP protocols %v and %v)",
			pc.Protocol,
			other.Protocol,
		)
	}

	// same protocols

	if pc.isTCPOrUDP() {
		portSet := pc.Ports.Intersect(*other.Ports)
		return NewProtocolContentWithPorts(pc.Protocol, &portSet), nil
	}

	if pc.isICMPv4OrICMPv6() {
		icmpSet := pc.ICMP.Intersect(*other.ICMP)
		return NewProtocolContentWithICMP(pc.Protocol, &icmpSet), nil
	}

	// custom Protocol

	return NewProtocolContentForCustomProtocol(pc.Protocol), nil
}

func (pc ProtocolContent) merge(other ProtocolContent) (ProtocolContent, error) {
	if other.noContent() {
		return pc, nil
	}

	if pc.allProtocols() || other.allProtocols() {
		return NewProtocolContentForAllTraffic(), nil
	}

	if pc.sameProtocolAs(other) == false {
		return ProtocolContent{}, fmt.Errorf(
			"cannot merge with different protocols (IP protocols %v and %v)",
			pc.Protocol,
			other.Protocol,
		)
	}

	// same protocols

	if pc.isTCPOrUDP() {
		portSet := pc.Ports.Merge(*other.Ports)
		return NewProtocolContentWithPorts(pc.Protocol, &portSet), nil
	}

	if pc.isICMPv4OrICMPv6() {
		icmpSet := pc.ICMP.Merge(*other.ICMP)
		return NewProtocolContentWithICMP(pc.Protocol, &icmpSet), nil
	}

	// custom Protocol

	return NewProtocolContentForCustomProtocol(pc.Protocol), nil
}

func (pc ProtocolContent) subtract(other ProtocolContent) (ProtocolContent, error) {
	if pc.noContent() {
		return pc, nil
	}

	if other.noContent() {
		return pc, nil
	}

	if other.allProtocols() {
		return NewProtocolContentForNoContent(), nil
	}

	// TODO: Handle subtracting one protocol from "all protocols"

	if pc.sameProtocolAs(other) == false {
		return ProtocolContent{}, fmt.Errorf(
			"cannot subtract with different protocols (IP protocols %v and %v)",
			pc.Protocol,
			other.Protocol,
		)
	}

	// same protocols

	if pc.isTCPOrUDP() {
		portSet := pc.Ports.Subtract(*other.Ports)
		return NewProtocolContentWithPorts(pc.Protocol, &portSet), nil
	}

	if pc.isICMPv4OrICMPv6() {
		icmpSet := pc.ICMP.Subtract(*other.ICMP)
		return NewProtocolContentWithICMP(pc.Protocol, &icmpSet), nil
	}

	// custom Protocol

	return NewProtocolContentForNoContent(), nil
}

func (pc ProtocolContent) sameProtocolAs(other ProtocolContent) bool {
	return pc.Protocol == other.Protocol
}

func (pc ProtocolContent) getProtocolName() string {
	switch pc.Protocol {
	case ProtocolAll:
		return ProtocolNameAll
	case ProtocolICMPv4:
		return ProtocolNameICMP
	case ProtocolTCP:
		return ProtocolNameTCP
	case ProtocolUDP:
		return ProtocolNameUDP
	case ProtocolICMPv6:
		return ProtocolNameICMPv6
	default:
		return string(pc.Protocol)
	}
}

func (pc ProtocolContent) usesNamedProtocol() bool {
	name := pc.getProtocolName()
	return strings.EqualFold(name, ProtocolNameTCP) ||
		strings.EqualFold(name, ProtocolNameUDP) ||
		strings.EqualFold(name, ProtocolNameICMP) ||
		strings.EqualFold(name, ProtocolNameICMPv6)
}
