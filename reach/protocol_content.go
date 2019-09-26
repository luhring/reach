package reach

import (
	"fmt"
	"strings"

	"github.com/luhring/reach/reach/set"
)

type ProtocolContent struct {
	Protocol Protocol
	Ports    *set.PortSet
	ICMP     *set.ICMPSet
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

func (p ProtocolContent) allProtocols() bool {
	return p.Protocol == ProtocolAll
}

func (p ProtocolContent) noContent() bool {
	return p.Protocol == ProtocolNone
}

func (p ProtocolContent) isTCPOrUDP() bool {
	return p.Protocol == ProtocolTCP || p.Protocol == ProtocolUDP
}

func (p ProtocolContent) isICMPv4OrICMPv6() bool {
	return p.Protocol == ProtocolICMPv4 || p.Protocol == ProtocolICMPv6
}

func (p ProtocolContent) intersect(other ProtocolContent) (ProtocolContent, error) {
	if p.noContent() || other.noContent() {
		return NewProtocolContentForNoContent(), nil
	}

	if p.allProtocols() {
		return other, nil
	}

	if other.allProtocols() {
		return p, nil
	}

	if p.sameProtocolAs(other) == false {
		return ProtocolContent{}, fmt.Errorf(
			"cannot intersect with different protocols (IP protocols %v and %v)",
			p.Protocol,
			other.Protocol,
		)
	}

	// same protocols

	if p.isTCPOrUDP() {
		portSet := p.Ports.Intersect(*other.Ports)
		return NewProtocolContentWithPorts(p.Protocol, &portSet), nil
	}

	if p.isICMPv4OrICMPv6() {
		icmpSet := p.ICMP.Intersect(*other.ICMP)
		return NewProtocolContentWithICMP(p.Protocol, &icmpSet), nil
	}

	// custom Protocol

	return NewProtocolContentForCustomProtocol(p.Protocol), nil
}

func (p ProtocolContent) merge(other ProtocolContent) (ProtocolContent, error) {
	if other.noContent() {
		return p, nil
	}

	if p.allProtocols() || other.allProtocols() {
		return NewProtocolContentForAllTraffic(), nil
	}

	if p.sameProtocolAs(other) == false {
		return ProtocolContent{}, fmt.Errorf(
			"cannot merge with different protocols (IP protocols %v and %v)",
			p.Protocol,
			other.Protocol,
		)
	}

	// same protocols

	if p.isTCPOrUDP() {
		portSet := p.Ports.Merge(*other.Ports)
		return NewProtocolContentWithPorts(p.Protocol, &portSet), nil
	}

	if p.isICMPv4OrICMPv6() {
		icmpSet := p.ICMP.Merge(*other.ICMP)
		return NewProtocolContentWithICMP(p.Protocol, &icmpSet), nil
	}

	// custom Protocol

	return NewProtocolContentForCustomProtocol(p.Protocol), nil
}

func (p ProtocolContent) subtract(other ProtocolContent) (ProtocolContent, error) {
	if p.noContent() {
		return p, nil
	}

	if other.noContent() {
		return p, nil
	}

	if other.allProtocols() {
		return NewProtocolContentForNoContent(), nil
	}

	// TODO: Handle subtracting one protocol from "all protocols"

	if p.sameProtocolAs(other) == false {
		return ProtocolContent{}, fmt.Errorf(
			"cannot subtract with different protocols (IP protocols %v and %v)",
			p.Protocol,
			other.Protocol,
		)
	}

	// same protocols

	if p.isTCPOrUDP() {
		portSet := p.Ports.Subtract(*other.Ports)
		return NewProtocolContentWithPorts(p.Protocol, &portSet), nil
	}

	if p.isICMPv4OrICMPv6() {
		icmpSet := p.ICMP.Subtract(*other.ICMP)
		return NewProtocolContentWithICMP(p.Protocol, &icmpSet), nil
	}

	// custom Protocol

	return NewProtocolContentForNoContent(), nil
}

func (p ProtocolContent) sameProtocolAs(other ProtocolContent) bool {
	return p.Protocol == other.Protocol
}

func (p ProtocolContent) getProtocolName() string {
	switch p.Protocol {
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
		return string(p.Protocol)
	}
}

func (p ProtocolContent) usesNamedProtocol() bool {
	name := p.getProtocolName()
	return strings.EqualFold(name, ProtocolNameTCP) ||
		strings.EqualFold(name, ProtocolNameUDP) ||
		strings.EqualFold(name, ProtocolNameICMP) ||
		strings.EqualFold(name, ProtocolNameICMPv6)
}
