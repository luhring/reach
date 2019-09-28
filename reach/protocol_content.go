package reach

import (
	"fmt"
	"log"
	"strings"

	"github.com/luhring/reach/reach/set"
)

type ProtocolContent struct {
	Protocol                 Protocol
	Ports                    *set.PortSet `json:"Ports,omitempty"`
	ICMP                     *set.ICMPSet `json:"ICMP,omitempty"`
	CustomProtocolHasContent *bool        `json:"CustomProtocolHasContent,omitempty"`
}

func NewProtocolContent(protocol Protocol, ports *set.PortSet, icmp *set.ICMPSet, customProtocolHasContent *bool) ProtocolContent {
	if protocol < 0 {
		log.Panicf("unexpected protocol value: %v", protocol) // TODO: Handle error better
	}

	return ProtocolContent{
		protocol,
		ports,
		icmp,
		customProtocolHasContent,
	}
}

func NewProtocolContentWithPorts(protocol Protocol, ports *set.PortSet) ProtocolContent {
	return NewProtocolContent(protocol, ports, nil, nil)
}

func NewProtocolContentWithPortsEmpty(protocol Protocol) ProtocolContent {
	ports := set.NewEmptyPortSet()
	return NewProtocolContentWithPorts(protocol, &ports)
}

func NewProtocolContentWithPortsFull(protocol Protocol) ProtocolContent {
	ports := set.NewFullPortSet()
	return NewProtocolContentWithPorts(protocol, &ports)
}

func NewProtocolContentWithICMP(protocol Protocol, icmp *set.ICMPSet) ProtocolContent {
	return NewProtocolContent(protocol, nil, icmp, nil)
}

func NewProtocolContentWithICMPEmpty(protocol Protocol) ProtocolContent {
	icmp := set.NewEmptyICMPSet()
	return NewProtocolContentWithICMP(protocol, &icmp)
}

func NewProtocolContentWithICMPFull(protocol Protocol) ProtocolContent {
	icmp := set.NewFullICMPSet()
	return NewProtocolContentWithICMP(protocol, &icmp)
}

func NewProtocolContentForCustomProtocol(protocol Protocol, hasContent bool) ProtocolContent {
	return NewProtocolContent(protocol, nil, nil, &hasContent)
}

func NewProtocolContentForCustomProtocolEmpty(protocol Protocol) ProtocolContent {
	hasContent := false
	return NewProtocolContent(protocol, nil, nil, &hasContent)
}

func NewProtocolContentForCustomProtocolFull(protocol Protocol) ProtocolContent {
	hasContent := true
	return NewProtocolContent(protocol, nil, nil, &hasContent)
}

func (pc ProtocolContent) Empty() bool {
	if pc.isTCPOrUDP() {
		return pc.Ports.Empty()
	} else if pc.isICMPv4OrICMPv6() {
		return pc.ICMP.Empty()
	} else {
		return !*pc.CustomProtocolHasContent
	}
}

func (pc ProtocolContent) String() string {
	output := ""

	if !pc.Empty() {
		output += ProtocolName(pc.Protocol)
		if pc.isTCPOrUDP() {
			output += fmt.Sprintf(": %s", pc.Ports.String())
		} else if pc.isICMPv4OrICMPv6() {
			output += fmt.Sprintf(": %s", pc.ICMP.String())
		}
	}

	return output
}

func (pc ProtocolContent) isTCPOrUDP() bool {
	return pc.Protocol == ProtocolTCP || pc.Protocol == ProtocolUDP
}

func (pc ProtocolContent) isICMPv4OrICMPv6() bool {
	return pc.Protocol == ProtocolICMPv4 || pc.Protocol == ProtocolICMPv6
}

func (pc ProtocolContent) intersect(other ProtocolContent) (ProtocolContent, error) {
	if !pc.sameProtocolAs(other) {
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

	if *pc.CustomProtocolHasContent && *other.CustomProtocolHasContent {
		return NewProtocolContentForCustomProtocolFull(pc.Protocol), nil
	}

	return NewProtocolContentForCustomProtocolEmpty(pc.Protocol), nil
}

func (pc ProtocolContent) merge(other ProtocolContent) (ProtocolContent, error) {
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

	if *pc.CustomProtocolHasContent || *other.CustomProtocolHasContent {
		return NewProtocolContentForCustomProtocol(pc.Protocol, true), nil
	}

	return NewProtocolContentForCustomProtocol(pc.Protocol, false), nil
}

func (pc ProtocolContent) subtract(other ProtocolContent) (ProtocolContent, error) {
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

	if *other.CustomProtocolHasContent || false == *pc.CustomProtocolHasContent {
		return NewProtocolContentForCustomProtocol(pc.Protocol, false), nil
	}

	return NewProtocolContentForCustomProtocol(pc.Protocol, true), nil
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
