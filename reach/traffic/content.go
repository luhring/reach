package traffic

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

// Content defines a set of network traffic across potentially multiple IP protocols.
type Content struct {
	indicator trafficContentIndicator
	protocols map[Protocol]*ProtocolContent
}

// NewContent returns a fresh, unset Content.
func NewContent() Content {
	return Content{
		indicator: trafficContentIndicatorUnset,
		protocols: nil,
	}
}

// All creates a new Content that represents the set of all expressible network traffic across all protocols.
func All() Content {
	return Content{
		indicator: trafficContentIndicatorAll,
	}
}

// None creates a new Content that represents a set of no network traffic.
func None() Content {
	return Content{
		indicator: trafficContentIndicatorNone,
	}
}

// ForPorts creates a new Content for a ports-oriented IP protocol, i.e. TCP or UDP.
func ForPorts(protocol Protocol, ports set.PortSet) Content {
	protocols := make(map[Protocol]*ProtocolContent)
	content := newProtocolContentWithPorts(protocol, ports)
	protocols[protocol] = &content

	return Content{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// ForICMP creates a new Content for either ICMPv4 or ICMPv6 traffic.
func ForICMP(protocol Protocol, icmp set.ICMPSet) Content {
	protocols := make(map[Protocol]*ProtocolContent)
	content := newProtocolContentWithICMP(protocol, &icmp)
	protocols[protocol] = &content

	return Content{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// ForCustomProtocol creates a new Content for a specified, custom IP protocol. The resulting Content will express either all traffic for that protocol or no traffic for that protocol, depending on the `hasContent` parameter.
func ForCustomProtocol(protocol Protocol, hasContent bool) Content {
	protocols := make(map[Protocol]*ProtocolContent)
	content := newProtocolContentForCustomProtocol(protocol, hasContent)
	protocols[protocol] = &content

	return Content{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

// Merge creates a new Content by merging any number of input TrafficContents.
func Merge(contents []Content) (Content, error) {
	result := NewContent()

	for _, trafficContent := range contents {
		if result.All() {
			return result, nil
		}

		mergedTrafficContent, err := result.Merge(trafficContent)
		if err != nil {
			return Content{}, err
		}

		result = mergedTrafficContent
	}

	return result, nil
}

// Intersect creates a new Content by intersecting any number of input TrafficContents.
func Intersect(contents []Content) Content {
	var result Content

	for i, trafficContent := range contents {
		if i == 0 {
			result = trafficContent
		} else {
			result = result.Intersect(trafficContent)

			if result.None() {
				return result
			}
		}
	}

	return result
}

// MergeTraffic returns the result of merging all input traffic contents. (MergeTraffic is a shortcut for Intersect.)
func MergeTraffic(tcs ...Content) Content {
	return Intersect(tcs)
}

// Merge performs a set merge operation on two TrafficContents.
func (c *Content) Merge(other Content) (Content, error) {
	if c.All() || other.All() {
		return All(), nil
	}

	if c.None() && other.None() {
		return None(), nil
	}

	result := NewContent()

	if !c.None() {
		for p := range c.protocols {
			mergedProtocolContent, err := result.Protocol(p).merge(c.Protocol(p))
			if err != nil {
				return Content{}, err
			}

			result.setProtocolContent(p, mergedProtocolContent)
		}
	}

	if !other.None() {
		for p := range other.protocols {
			mergedProtocolContent, err := result.Protocol(p).merge(other.Protocol(p))
			if err != nil {
				return Content{}, err
			}

			result.setProtocolContent(p, mergedProtocolContent)
		}
	}

	return result, nil
}

// Intersect performs a set intersection operation on two TrafficContents.
func (c *Content) Intersect(other Content) Content {
	if c.None() || other.None() {
		return None()
	}

	if c.All() && other.All() {
		return All()
	}

	protocolsToProcess := make(map[Protocol]bool)

	if !c.All() {
		for p := range c.protocols {
			protocolsToProcess[p] = true
		}
	}

	if !other.All() {
		for p := range other.protocols {
			protocolsToProcess[p] = true
		}
	}

	result := NewContent()

	for p, shouldProcess := range protocolsToProcess {
		if shouldProcess && !c.Protocol(p).Empty() && !other.Protocol(p).Empty() {
			intersection := c.Protocol(p).intersect(other.Protocol(p))
			result.setProtocolContent(p, intersection)
		}
	}

	return result
}

// Subtract performs a set subtraction (self - other) on two TrafficContents.
func (c *Content) Subtract(other Content) (Content, error) {
	if c.None() || other.All() {
		return None(), nil
	}

	if other.None() {
		return *c, nil
	}

	result := NewContent()

	for p, pc := range c.protocols {
		pcDifference, err := pc.subtract(other.Protocol(p))
		if err != nil {
			return Content{}, fmt.Errorf("unable to subtract traffic content: %v", err)
		}

		result.setProtocolContent(p, pcDifference)
	}

	return result, nil
}

// MarshalJSON returns the JSON representation of the Content.
func (c Content) MarshalJSON() ([]byte, error) {
	if c.None() {
		return json.Marshal("[no traffic]")
	}

	if c.All() {
		return json.Marshal("[all traffic]")
	}

	result := make(map[string][]string)

	for protocol, content := range c.protocols {
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

// String returns the string representation of the Content.
func (c Content) String() string {
	if c.All() {
		return allTrafficString + "\n"
	}

	if c.None() {
		return noTrafficString + "\n"
	}

	var tcpLines, udpLines, icmpv4Lines, icmpv6Lines []string
	var customProtocolContents []*ProtocolContent
	var customOutputItems []string
	var outputItems []string

	for _, content := range c.protocols {
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

// StringWithSymbols returns the string representation of the Content, with the added feature of pre-pending each output line with a symbol, intended for display to the user.
func (c Content) StringWithSymbols() string {
	if c.All() {
		return "✓ " + allTrafficString + "\n"
	}

	if c.None() {
		return noTrafficString + "\n"
	}

	var tcpLines, udpLines, icmpv4Lines, icmpv6Lines []string
	var customProtocolContents []*ProtocolContent
	var customOutputItems []string
	var outputItems []string

	for _, content := range c.protocols {
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

// ColorString returns the string representation of the Content, where the positive traffic findings are displayed as green, and the absence of traffic is displayed as red.
func (c Content) ColorString() string {
	if c.None() {
		return ansi.Color(c.String(), "red+b")
	}
	return ansi.Color(c.String(), "green+b")
}

// ColorStringWithSymbols returns the colored version of the output from StringWithSymbols().
func (c Content) ColorStringWithSymbols() string {
	if c.None() {
		return ansi.Color(c.StringWithSymbols(), "red+b")
	}
	return ansi.Color(c.StringWithSymbols(), "green+b")
}

// All returns a boolean indicating whether or not the Content represents all network traffic.
func (c Content) All() bool {
	return c.indicator == trafficContentIndicatorAll
}

// None returns a boolean indicating whether or not the Content represents no network traffic.
func (c Content) None() bool {
	return c.indicator == trafficContentIndicatorNone || (c.indicator == trafficContentIndicatorUnset && len(c.protocols) == 0)
}

// Protocols returns a slice of the IP protocols described by the traffic content.
func (c Content) Protocols() []Protocol {
	if c.protocols == nil {
		return nil
	}

	var result []Protocol

	for protocol := range c.protocols {
		result = append(result, protocol)
	}

	return result
}

func (c *Content) setProtocolContent(p Protocol, content ProtocolContent) {
	c.indicator = trafficContentIndicatorUnset

	if c.protocols == nil {
		c.protocols = make(map[Protocol]*ProtocolContent)
	}

	c.protocols[p] = &content
}

// Protocol returns the protocol-specific content within the Content for the specified IP protocol.
func (c Content) Protocol(p Protocol) ProtocolContent {
	content := c.protocols[p]

	if content == nil {
		if p.UsesPorts() {
			if c.All() {
				return newProtocolContentWithPortsFull(p)
			}
			return newProtocolContentWithPortsEmpty(p)
		}

		if p.UsesICMPTypeCodes() {
			if c.All() {
				return newProtocolContentWithICMPFull(p)
			}
			return newProtocolContentWithICMPEmpty(p)
		}

		// custom protocol

		if c.All() {
			return newProtocolContentForCustomProtocolFull(p)
		}
		return newProtocolContentForCustomProtocolEmpty(p)
	}

	return *content
}

// HasProtocol returns a bool to indicate whether the specified protocol exists within the traffic content.
func (c Content) HasProtocol(p Protocol) bool {
	return c.Protocol(p).Empty() == false
}
