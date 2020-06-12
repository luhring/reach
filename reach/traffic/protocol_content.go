package traffic

import (
	"fmt"
	"log"
	"strings"

	"github.com/luhring/reach/reach/set"
)

// ProtocolContent specifies a set of network traffic for a single, specified IP protocol.
type ProtocolContent struct {
	Protocol                 Protocol
	Ports                    *set.PortSet `json:"Ports,omitempty"`
	ICMP                     *set.ICMPSet `json:"ICMP,omitempty"`
	CustomProtocolHasContent *bool        `json:"CustomProtocolHasContent,omitempty"`
}

func newProtocolContent(protocol Protocol, ports *set.PortSet, icmp *set.ICMPSet, customProtocolHasContent *bool) ProtocolContent {
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

func newProtocolContentWithPorts(protocol Protocol, ports set.PortSet) ProtocolContent {
	return newProtocolContent(protocol, &ports, nil, nil)
}

func newProtocolContentWithPortsEmpty(protocol Protocol) ProtocolContent {
	ports := set.NewEmptyPortSet()
	return newProtocolContentWithPorts(protocol, ports)
}

func newProtocolContentWithPortsFull(protocol Protocol) ProtocolContent {
	ports := set.NewFullPortSet()
	return newProtocolContentWithPorts(protocol, ports)
}

func newProtocolContentWithICMP(protocol Protocol, icmp *set.ICMPSet) ProtocolContent {
	return newProtocolContent(protocol, nil, icmp, nil)
}

func newProtocolContentWithICMPEmpty(protocol Protocol) ProtocolContent {
	icmp := set.NewEmptyICMPSet()
	return newProtocolContentWithICMP(protocol, &icmp)
}

func newProtocolContentWithICMPFull(protocol Protocol) ProtocolContent {
	icmp := set.NewFullICMPSet()
	return newProtocolContentWithICMP(protocol, &icmp)
}

func newProtocolContentForCustomProtocol(protocol Protocol, hasContent bool) ProtocolContent {
	return newProtocolContent(protocol, nil, nil, &hasContent)
}

func newProtocolContentForCustomProtocolEmpty(protocol Protocol) ProtocolContent {
	hasContent := false
	return newProtocolContent(protocol, nil, nil, &hasContent)
}

func newProtocolContentForCustomProtocolFull(protocol Protocol) ProtocolContent {
	hasContent := true
	return newProtocolContent(protocol, nil, nil, &hasContent)
}

// Empty returns a bool indicating whether this ProtocolContent represents no traffic for this protocol.
func (pc ProtocolContent) Empty() bool {
	if pc.isTCPOrUDP() {
		return pc.Ports == nil || pc.Ports.Empty()
	} else if pc.isICMPv4OrICMPv6() {
		return pc.ICMP == nil || pc.ICMP.Empty()
	} else {
		return !*pc.CustomProtocolHasContent
	}
}

// Complete returns a bool indicating whether this ProtocolContent represents all traffic for this protocol.
func (pc ProtocolContent) Complete() bool {
	if pc.isTCPOrUDP() {
		return pc.Ports.Complete()
	} else if pc.isICMPv4OrICMPv6() {
		return pc.ICMP.Complete()
	} else {
		return *pc.CustomProtocolHasContent
	}
}

// String returns the string representation of the protocol content.
func (pc ProtocolContent) String() string {
	protocolName := ProtocolName(pc.Protocol)

	if !pc.Empty() {
		if pc.isTCPOrUDP() {
			return fmt.Sprintf("%s %s", protocolName, pc.Ports.String())
		} else if pc.isICMPv4OrICMPv6() {
			if pc.Protocol == ProtocolICMPv6 {
				return fmt.Sprintf("%s", pc.ICMP.StringV6())
			}
			return fmt.Sprintf("%s", pc.ICMP.StringV4())
		}
		return fmt.Sprintf("%s (all traffic)", protocolName)
	}
	return fmt.Sprintf("%s (no traffic)", protocolName)
}

func (pc ProtocolContent) lines() []string {
	protocolName := ProtocolName(pc.Protocol)

	if !pc.Empty() {
		if pc.isTCPOrUDP() {

			var lines []string

			for _, rangeString := range pc.Ports.RangeStrings() {
				lines = append(lines, fmt.Sprintf("%s %s", protocolName, rangeString))
			}
			return lines
		} else if pc.isICMPv4OrICMPv6() {
			if pc.Protocol == ProtocolICMPv6 {
				return pc.ICMP.RangeStringsV6()
			}
			return pc.ICMP.RangeStringsV4()
		} else {
			return []string{fmt.Sprintf("%s (all traffic)", protocolName)}
		}
	}

	return []string{fmt.Sprintf("%s (no traffic)", protocolName)}
}

func (pc ProtocolContent) isTCPOrUDP() bool {
	return pc.Protocol == ProtocolTCP || pc.Protocol == ProtocolUDP
}

func (pc ProtocolContent) isICMPv4OrICMPv6() bool {
	return pc.Protocol == ProtocolICMPv4 || pc.Protocol == ProtocolICMPv6
}

func (pc ProtocolContent) intersect(other ProtocolContent) ProtocolContent {
	if !pc.sameProtocolAs(other) {
		return newProtocolContentForCustomProtocolEmpty(pc.Protocol)
	}

	// same protocols

	if pc.isTCPOrUDP() {
		portSet := pc.Ports.Intersect(*other.Ports)
		return newProtocolContentWithPorts(pc.Protocol, portSet)
	}

	if pc.isICMPv4OrICMPv6() {
		icmpSet := pc.ICMP.Intersect(*other.ICMP)
		return newProtocolContentWithICMP(pc.Protocol, &icmpSet)
	}

	// custom Protocol

	if *pc.CustomProtocolHasContent && *other.CustomProtocolHasContent {
		return newProtocolContentForCustomProtocolFull(pc.Protocol)
	}

	return newProtocolContentForCustomProtocolEmpty(pc.Protocol)
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
		return newProtocolContentWithPorts(pc.Protocol, portSet), nil
	}

	if pc.isICMPv4OrICMPv6() {
		icmpSet := pc.ICMP.Merge(*other.ICMP)
		return newProtocolContentWithICMP(pc.Protocol, &icmpSet), nil
	}

	// custom Protocol

	if *pc.CustomProtocolHasContent || *other.CustomProtocolHasContent {
		return newProtocolContentForCustomProtocol(pc.Protocol, true), nil
	}

	return newProtocolContentForCustomProtocol(pc.Protocol, false), nil
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
		return newProtocolContentWithPorts(pc.Protocol, portSet), nil
	}

	if pc.isICMPv4OrICMPv6() {
		icmpSet := pc.ICMP.Subtract(*other.ICMP)
		return newProtocolContentWithICMP(pc.Protocol, &icmpSet), nil
	}

	// custom Protocol

	if *other.CustomProtocolHasContent || false == *pc.CustomProtocolHasContent {
		return newProtocolContentForCustomProtocol(pc.Protocol, false), nil
	}

	return newProtocolContentForCustomProtocol(pc.Protocol, true), nil
}

func (pc ProtocolContent) sameProtocolAs(other ProtocolContent) bool {
	return pc.Protocol == other.Protocol
}

func (pc ProtocolContent) getProtocolName() string {
	switch pc.Protocol {
	case ProtocolAll:
		return ProtocolNameAll
	case ProtocolICMPv4:
		return ProtocolNameICMPv4
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
		strings.EqualFold(name, ProtocolNameICMPv4) ||
		strings.EqualFold(name, ProtocolNameICMPv6)
}
