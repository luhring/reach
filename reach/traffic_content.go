package reach

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/mgutz/ansi"

	"github.com/luhring/reach/reach/helper"
	"github.com/luhring/reach/reach/set"
)

const (
	trafficContentIndicatorUnset trafficContentIndicator = iota
	trafficContentIndicatorAll
	trafficContentIndicatorNone
	allTrafficString = "all traffic"
	noTrafficString  = "(none)"
)

type trafficContentIndicator int

type TrafficContent struct {
	indicator trafficContentIndicator
	protocols map[Protocol]*ProtocolContent
}

func NewTrafficContent() TrafficContent {
	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: nil,
	}
}

func NewTrafficContentForAllTraffic() TrafficContent {
	return TrafficContent{
		indicator: trafficContentIndicatorAll,
	}
}

func NewTrafficContentForNoTraffic() TrafficContent {
	return TrafficContent{
		indicator: trafficContentIndicatorNone,
	}
}

func NewTrafficContentForPorts(protocol Protocol, ports set.PortSet) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	content := NewProtocolContentWithPorts(protocol, &ports)
	protocols[protocol] = &content

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentForICMP(protocol Protocol, icmp set.ICMPSet) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	content := NewProtocolContentWithICMP(protocol, &icmp)
	protocols[protocol] = &content

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentForCustomProtocol(protocol Protocol, hasContent bool) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	content := NewProtocolContentForCustomProtocol(protocol, hasContent)
	protocols[protocol] = &content

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentFromMergingMultiple(contents []TrafficContent) (TrafficContent, error) {
	result := NewTrafficContent()

	for _, trafficContent := range contents {
		if result.All() {
			return result, nil
		}

		mergedTrafficContent, err := result.Merge(trafficContent)
		if err != nil {
			return TrafficContent{}, err
		}

		result = mergedTrafficContent
	}

	return result, nil
}

func NewTrafficContentFromIntersectingMultiple(contents []TrafficContent) (TrafficContent, error) {
	var result TrafficContent

	for i, trafficContent := range contents {
		if i == 0 {
			result = trafficContent
		} else {
			intersection, err := result.Intersect(trafficContent)
			if err != nil {
				return TrafficContent{}, err
			}

			result = intersection

			if result.None() {
				return result, nil
			}
		}
	}

	return result, nil
}

func TrafficContentsFromFactors(factors []Factor) []TrafficContent {
	var result []TrafficContent

	for _, factor := range factors {
		result = append(result, factor.Traffic)
	}

	return result
}

func (tc *TrafficContent) Merge(other TrafficContent) (TrafficContent, error) {
	if tc.All() || other.All() {
		return NewTrafficContentForAllTraffic(), nil
	}

	if tc.None() && other.None() {
		return NewTrafficContentForNoTraffic(), nil
	}

	result := NewTrafficContent()

	if !tc.None() {
		for p := range tc.protocols {
			mergedProtocolContent, err := result.Protocol(p).merge(tc.Protocol(p))
			if err != nil {
				return TrafficContent{}, err
			}

			result.SetProtocolContent(p, mergedProtocolContent)
		}
	}

	if !other.None() {
		for p := range other.protocols {
			mergedProtocolContent, err := result.Protocol(p).merge(other.Protocol(p))
			if err != nil {
				return TrafficContent{}, err
			}

			result.SetProtocolContent(p, mergedProtocolContent)
		}
	}

	return result, nil
}

func (tc *TrafficContent) Intersect(other TrafficContent) (TrafficContent, error) {
	if tc.None() || other.None() {
		return NewTrafficContentForNoTraffic(), nil
	}

	if tc.All() && other.All() {
		return NewTrafficContentForAllTraffic(), nil
	}

	protocolsToProcess := make(map[Protocol]bool)

	if !tc.All() {
		for p := range tc.protocols {
			protocolsToProcess[p] = true
		}
	}

	if !other.All() {
		for p := range other.protocols {
			protocolsToProcess[p] = true
		}
	}

	result := NewTrafficContent()

	for p, shouldProcess := range protocolsToProcess {
		if shouldProcess && !tc.Protocol(p).Empty() && !other.Protocol(p).Empty() {
			intersection, err := tc.Protocol(p).intersect(other.Protocol(p))
			if err != nil {
				return TrafficContent{}, err
			}

			result.SetProtocolContent(p, intersection)
		}
	}

	return result, nil
}

func (tc TrafficContent) MarshalJSON() ([]byte, error) {
	if tc.None() {
		return json.Marshal("[no traffic]")
	}

	if tc.All() {
		return json.Marshal("[all traffic]")
	}

	result := make(map[string][]string)

	for protocol, content := range tc.protocols {
		key := ProtocolName(protocol)
		if protocol.UsesPorts() {
			result[key] = content.Ports.RangeStrings()
		} else if protocol.UsesICMPTypeCodes() {
			if protocol == ProtocolICMPv6 {
				result[key] = content.ICMP.RangeStringsV6()
			} else {
				result[key] = content.ICMP.RangeStringsV4()
			}
		} else {
			if content.CustomProtocolHasContent != nil && *content.CustomProtocolHasContent {
				result[key] = []string{"[all traffic]"}
			} else {
				result[key] = []string{"[no traffic]"}
			}
		}
	}

	return json.Marshal(result)
}

func (tc TrafficContent) String() string {
	if tc.All() {
		return allTrafficString + "\n"
	}

	if tc.None() {
		return noTrafficString + "\n"
	}

	var tcpLines, udpLines, icmpv4Lines, icmpv6Lines []string
	var customProtocolContents []*ProtocolContent
	var customOutputItems []string
	var outputItems []string

	for _, content := range tc.protocols {
		switch content.Protocol {
		case ProtocolTCP:
			tcpLines = append(tcpLines, content.Lines()...)
		case ProtocolUDP:
			udpLines = append(udpLines, content.Lines()...)
		case ProtocolICMPv4:
			icmpv4Lines = append(icmpv4Lines, content.Lines()...)
		case ProtocolICMPv6:
			icmpv6Lines = append(icmpv6Lines, content.Lines()...)
		default:
			customProtocolContents = append(customProtocolContents, content)
		}
	}
	sort.Slice(customProtocolContents, func(i, j int) bool {
		return customProtocolContents[i].Protocol < customProtocolContents[j].Protocol
	})

	for _, item := range customProtocolContents {
		customOutputItems = append(customOutputItems, item.String())
	}

	customOutput := strings.Join(customOutputItems, "\n")

	if len(tcpLines) > 0 {
		outputItems = append(outputItems, tcpLines...)
	}
	if len(udpLines) > 0 {
		outputItems = append(outputItems, udpLines...)
	}
	if len(icmpv4Lines) > 0 {
		outputItems = append(outputItems, icmpv4Lines...)
	}
	if len(icmpv6Lines) > 0 {
		outputItems = append(outputItems, icmpv6Lines...)
	}
	if customOutput != "" {
		outputItems = append(outputItems, customOutput)
	}

	return strings.Join(outputItems, "\n") + "\n"
}

func (tc TrafficContent) StringWithSymbols() string {
	if tc.All() {
		return "✓ " + allTrafficString + "\n"
	}

	if tc.None() {
		return noTrafficString + "\n"
	}

	var tcpLines, udpLines, icmpv4Lines, icmpv6Lines []string
	var customProtocolContents []*ProtocolContent
	var customOutputItems []string
	var outputItems []string

	for _, content := range tc.protocols {
		switch content.Protocol {
		case ProtocolTCP:
			tcpLines = append(tcpLines, content.Lines()...)
		case ProtocolUDP:
			udpLines = append(udpLines, content.Lines()...)
		case ProtocolICMPv4:
			icmpv4Lines = append(icmpv4Lines, content.Lines()...)
		case ProtocolICMPv6:
			icmpv6Lines = append(icmpv6Lines, content.Lines()...)
		default:
			customProtocolContents = append(customProtocolContents, content)
		}
	}
	sort.Slice(customProtocolContents, func(i, j int) bool {
		return customProtocolContents[i].Protocol < customProtocolContents[j].Protocol
	})

	for _, item := range customProtocolContents {
		customOutputItems = append(customOutputItems, "✓ "+item.String())
	}

	customOutput := strings.Join(customOutputItems, "\n")

	if len(tcpLines) > 0 {
		outputItems = append(outputItems, helper.PrefixLines(tcpLines, "✓ "))
	}
	if len(udpLines) > 0 {
		outputItems = append(outputItems, helper.PrefixLines(udpLines, "✓ "))
	}
	if len(icmpv4Lines) > 0 {
		outputItems = append(outputItems, helper.PrefixLines(icmpv4Lines, "✓ "))
	}
	if len(icmpv6Lines) > 0 {
		outputItems = append(outputItems, helper.PrefixLines(icmpv6Lines, "✓ "))
	}
	if customOutput != "" {
		outputItems = append(outputItems, customOutput)
	}

	return strings.Join(outputItems, "\n") + "\n"
}

func (tc TrafficContent) ColorString() string {
	if tc.None() {
		return ansi.Color(tc.String(), "red+b")
	} else {
		return ansi.Color(tc.String(), "green+b")
	}
}

func (tc TrafficContent) ColorStringWithSymbols() string {
	if tc.None() {
		return ansi.Color(tc.StringWithSymbols(), "red+b")
	} else {
		return ansi.Color(tc.StringWithSymbols(), "green+b")
	}
}

func (tc TrafficContent) All() bool {
	return tc.indicator == trafficContentIndicatorAll
}

func (tc TrafficContent) None() bool {
	return tc.indicator == trafficContentIndicatorNone || (tc.indicator == trafficContentIndicatorUnset && len(tc.protocols) == 0)
}

func (tc *TrafficContent) SetAll() {
	tc.indicator = trafficContentIndicatorAll
	tc.protocols = nil
}

func (tc *TrafficContent) SetNone() {
	tc.indicator = trafficContentIndicatorNone
	tc.protocols = nil
}

func (tc *TrafficContent) SetProtocolContent(p Protocol, content ProtocolContent) {
	tc.indicator = trafficContentIndicatorUnset

	if tc.protocols == nil {
		tc.protocols = make(map[Protocol]*ProtocolContent)
	}

	tc.protocols[p] = &content
}

func (tc TrafficContent) Protocol(p Protocol) ProtocolContent {
	content := tc.protocols[p]

	if content == nil {
		if p.UsesPorts() {
			if tc.All() {
				return NewProtocolContentWithPortsFull(p)
			}
			return NewProtocolContentWithPortsEmpty(p)
		}

		if p.UsesICMPTypeCodes() {
			if tc.All() {
				return NewProtocolContentWithICMPFull(p)
			}
			return NewProtocolContentWithICMPEmpty(p)
		}

		// custom protocol

		if tc.All() {
			return NewProtocolContentForCustomProtocolFull(p)
		}
		return NewProtocolContentForCustomProtocolEmpty(p)
	}

	return *content
}

func (tc TrafficContent) TCP() ProtocolContent {
	return tc.Protocol(ProtocolTCP)
}

func (tc TrafficContent) UDP() ProtocolContent {
	return tc.Protocol(ProtocolUDP)
}

func (tc TrafficContent) ICMPv4() ProtocolContent {
	return tc.Protocol(ProtocolICMPv4)
}

func (tc TrafficContent) ICMPv6() ProtocolContent {
	return tc.Protocol(ProtocolICMPv6)
}
