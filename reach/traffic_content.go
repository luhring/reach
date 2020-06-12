package reach

import (
	"encoding/json"
	"fmt"
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

// TrafficContent defines a set of network traffic across potentially multiple IP protocols.
type TrafficContent struct {
	indicator trafficContentIndicator
	protocols map[Protocol]*ProtocolContent
}

func newTrafficContent() TrafficContent {
	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: nil,
	}
}

// NewTrafficContentForAllTraffic creates a new TrafficContent that represents the set of all expressible network traffic across all protocols.
func NewTrafficContentForAllTraffic() TrafficContent {
	return TrafficContent{
		indicator: trafficContentIndicatorAll,
	}
}

// NewTrafficContentForNoTraffic creates a new TrafficContent that represents a set of no network traffic.
func NewTrafficContentForNoTraffic() TrafficContent {
	return TrafficContent{
		indicator: trafficContentIndicatorNone,
	}
}

// NewTrafficContentForPorts creates a new TrafficContent for a ports-oriented IP protocol, i.e. TCP or UDP.
func NewTrafficContentForPorts(protocol Protocol, ports set.PortSet) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	content := newProtocolContentWithPorts(protocol, &ports)
	protocols[protocol] = &content

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// NewTrafficContentForICMP creates a new TrafficContent for either ICMPv4 or ICMPv6 traffic.
func NewTrafficContentForICMP(protocol Protocol, icmp set.ICMPSet) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	content := newProtocolContentWithICMP(protocol, &icmp)
	protocols[protocol] = &content

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// NewTrafficContentForCustomProtocol creates a new TrafficContent for a specified, custom IP protocol. The resulting TrafficContent will express either all traffic for that protocol or no traffic for that protocol, depending on the `hasContent` parameter.
func NewTrafficContentForCustomProtocol(protocol Protocol, hasContent bool) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	content := newProtocolContentForCustomProtocol(protocol, hasContent)
	protocols[protocol] = &content

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// NewTrafficContentFromMergingMultiple creates a new TrafficContent by merging any number of input TrafficContents.
func NewTrafficContentFromMergingMultiple(contents []TrafficContent) (TrafficContent, error) {
	result := newTrafficContent()

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

// NewTrafficContentFromIntersectingMultiple creates a new TrafficContent by intersecting any number of input TrafficContents.
func NewTrafficContentFromIntersectingMultiple(contents []TrafficContent) TrafficContent {
	var result TrafficContent

	for i, trafficContent := range contents {
		if i == 0 {
			result = trafficContent
		} else {
			result := result.Intersect(trafficContent)

			if result.None() {
				return result
			}
		}
	}

	return result
}

// MergeTraffic returns the result of merging all input traffic contents. (MergeTraffic is a shortcut for NewTrafficContentFromIntersectingMultiple.)
func MergeTraffic(tcs ...TrafficContent) TrafficContent {
	return NewTrafficContentFromIntersectingMultiple(tcs)
}

// Merge performs a set merge operation on two TrafficContents.
func (tc *TrafficContent) Merge(other TrafficContent) (TrafficContent, error) {
	if tc.All() || other.All() {
		return NewTrafficContentForAllTraffic(), nil
	}

	if tc.None() && other.None() {
		return NewTrafficContentForNoTraffic(), nil
	}

	result := newTrafficContent()

	if !tc.None() {
		for p := range tc.protocols {
			mergedProtocolContent, err := result.Protocol(p).merge(tc.Protocol(p))
			if err != nil {
				return TrafficContent{}, err
			}

			result.setProtocolContent(p, mergedProtocolContent)
		}
	}

	if !other.None() {
		for p := range other.protocols {
			mergedProtocolContent, err := result.Protocol(p).merge(other.Protocol(p))
			if err != nil {
				return TrafficContent{}, err
			}

			result.setProtocolContent(p, mergedProtocolContent)
		}
	}

	return result, nil
}

// Intersect performs a set intersection operation on two TrafficContents.
func (tc *TrafficContent) Intersect(other TrafficContent) TrafficContent {
	if tc.None() || other.None() {
		return NewTrafficContentForNoTraffic()
	}

	if tc.All() && other.All() {
		return NewTrafficContentForAllTraffic()
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

	result := newTrafficContent()

	for p, shouldProcess := range protocolsToProcess {
		if shouldProcess && !tc.Protocol(p).Empty() && !other.Protocol(p).Empty() {
			intersection := tc.Protocol(p).intersect(other.Protocol(p))
			result.setProtocolContent(p, intersection)
		}
	}

	return result
}

// Subtract performs a set subtraction (self - other) on two TrafficContents.
func (tc *TrafficContent) Subtract(other TrafficContent) (TrafficContent, error) {
	if tc.None() || other.All() {
		return NewTrafficContentForNoTraffic(), nil
	}

	if other.None() {
		return *tc, nil
	}

	result := newTrafficContent()

	for p, pc := range tc.protocols {
		pcDifference, err := pc.subtract(other.Protocol(p))
		if err != nil {
			return TrafficContent{}, fmt.Errorf("unable to subtract traffic content: %v", err)
		}

		result.setProtocolContent(p, pcDifference)
	}

	return result, nil
}

// MarshalJSON returns the JSON representation of the TrafficContent.
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

// String returns the string representation of the TrafficContent.
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
			tcpLines = append(tcpLines, content.lines()...)
		case ProtocolUDP:
			udpLines = append(udpLines, content.lines()...)
		case ProtocolICMPv4:
			icmpv4Lines = append(icmpv4Lines, content.lines()...)
		case ProtocolICMPv6:
			icmpv6Lines = append(icmpv6Lines, content.lines()...)
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

// StringWithSymbols returns the string representation of the TrafficContent, with the added feature of pre-pending each output line with a symbol, intended for display to the user.
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
			tcpLines = append(tcpLines, content.lines()...)
		case ProtocolUDP:
			udpLines = append(udpLines, content.lines()...)
		case ProtocolICMPv4:
			icmpv4Lines = append(icmpv4Lines, content.lines()...)
		case ProtocolICMPv6:
			icmpv6Lines = append(icmpv6Lines, content.lines()...)
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

// ColorString returns the string representation of the TrafficContent, where the positive traffic findings are displayed as green, and the absence of traffic is displayed as red.
func (tc TrafficContent) ColorString() string {
	if tc.None() {
		return ansi.Color(tc.String(), "red+b")
	}
	return ansi.Color(tc.String(), "green+b")
}

// ColorStringWithSymbols returns the colored version of the output from StringWithSymbols().
func (tc TrafficContent) ColorStringWithSymbols() string {
	if tc.None() {
		return ansi.Color(tc.StringWithSymbols(), "red+b")
	}
	return ansi.Color(tc.StringWithSymbols(), "green+b")
}

// All returns a boolean indicating whether or not the TrafficContent represents all network traffic.
func (tc TrafficContent) All() bool {
	return tc.indicator == trafficContentIndicatorAll
}

// None returns a boolean indicating whether or not the TrafficContent represents no network traffic.
func (tc TrafficContent) None() bool {
	return tc.indicator == trafficContentIndicatorNone || (tc.indicator == trafficContentIndicatorUnset && len(tc.protocols) == 0)
}

// Protocols returns a slice of the IP protocols described by the traffic content.
func (tc TrafficContent) Protocols() []Protocol {
	if tc.protocols == nil {
		return nil
	}

	var result []Protocol

	for protocol := range tc.protocols {
		result = append(result, protocol)
	}

	return result
}

func (tc *TrafficContent) setProtocolContent(p Protocol, content ProtocolContent) {
	tc.indicator = trafficContentIndicatorUnset

	if tc.protocols == nil {
		tc.protocols = make(map[Protocol]*ProtocolContent)
	}

	tc.protocols[p] = &content
}

// Protocol returns the protocol-specific content within the TrafficContent for the specified IP protocol.
func (tc TrafficContent) Protocol(p Protocol) ProtocolContent {
	content := tc.protocols[p]

	if content == nil {
		if p.UsesPorts() {
			if tc.All() {
				return newProtocolContentWithPortsFull(p)
			}
			return newProtocolContentWithPortsEmpty(p)
		}

		if p.UsesICMPTypeCodes() {
			if tc.All() {
				return newProtocolContentWithICMPFull(p)
			}
			return newProtocolContentWithICMPEmpty(p)
		}

		// custom protocol

		if tc.All() {
			return newProtocolContentForCustomProtocolFull(p)
		}
		return newProtocolContentForCustomProtocolEmpty(p)
	}

	return *content
}

// HasProtocol returns a bool to indicate whether the specified protocol exists within the traffic content.
func (tc TrafficContent) HasProtocol(p Protocol) bool {
	return tc.Protocol(p).Empty() == false
}
