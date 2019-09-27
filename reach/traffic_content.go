package reach

import (
	"encoding/json"
	"fmt"

	"github.com/luhring/reach/reach/set"
)

const (
	trafficContentIndicatorUnset trafficContentIndicator = iota
	trafficContentIndicatorAll
	trafficContentIndicatorNone
)

type trafficContentIndicator int

type TrafficContent struct {
	indicator trafficContentIndicator
	protocols map[Protocol]*ProtocolContent
}

func NewTrafficContent() *TrafficContent {
	return &TrafficContent{
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
	protocols[protocol] = NewProtocolContentWithPorts(protocol, &ports)

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentForICMP(protocol Protocol, icmp set.ICMPSet) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	protocols[protocol] = NewProtocolContentWithICMP(protocol, &icmp)

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentForCustomProtocol(protocol Protocol, hasContent bool) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	protocols[protocol] = NewProtocolContentForCustomProtocol(protocol, hasContent)

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentFromMergingMultiple(contents []TrafficContent) (*TrafficContent, error) {
	result := NewTrafficContent()

	for _, tc := range contents {
		if result.All() {
			return result, nil
		}

		err := result.Merge(tc)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func NewTrafficContentFromIntersectingMultiple(contents []TrafficContent) (*TrafficContent, error) {
	var content TrafficContent

	for i, tc := range contents {
		if i == 0 {
			content = tc
		} else {
			err := content.Intersect(tc)
			if err != nil {
				return nil, err
			}

			if content.None() {
				return &content, nil
			}
		}

	}

	return &content, nil
}

func (tc *TrafficContent) Merge(other TrafficContent) error {
	if tc.All() || other.All() {
		tc.SetAll()
		return nil
	}

	if tc.None() && other.None() {
		tc.SetNone()
		return nil
	}

	for protocol, content := range other.protocols {
		mergedProtocolContent, err := tc.Protocol(protocol).merge(*content)
		if err != nil {
			return err
		}

		tc.SetProtocolContent(protocol, mergedProtocolContent)
	}

	return nil
}

func (tc *TrafficContent) Intersect(other TrafficContent) error {
	if tc.None() || other.None() {
		tc.SetNone()
		return nil
	}

	if tc.All() && other.All() {
		tc.SetAll()
		return nil
	}

	for protocol, content := range other.protocols {
		intersectedProtocolContent, err := tc.Protocol(protocol).intersect(*content)
		if err != nil {
			return err
		}

		tc.SetProtocolContent(protocol, intersectedProtocolContent)
	}

	return nil
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
			result[key] = content.ICMP.Types()
		} else {
			if content.CustomProtocolHasContent != nil && *content.CustomProtocolHasContent {
				result[key] = []string{"[all traffic for this protocol]"}
			} else {
				result[key] = []string{"[no traffic for this protocol]"}
			}
		}
	}

	return json.Marshal(result)
}

func (tc TrafficContent) String() string {
	if tc.All() {
		return "[all traffic]"
	}

	if tc.None() {
		return "[no traffic]"
	}

	var output string

	for p, content := range tc.protocols {
		output += ProtocolName(p)
		if p.UsesPorts() {
			output += fmt.Sprintf(": %s\n", content.Ports.String())
		} else if p.UsesICMPTypeCodes() {
			output += fmt.Sprintf(": %s\n", content.ICMP.String())
		}
	}

	return output
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

func (tc *TrafficContent) SetProtocolContent(p Protocol, content *ProtocolContent) {
	tc.indicator = trafficContentIndicatorUnset

	if tc.protocols == nil {
		tc.protocols = make(map[Protocol]*ProtocolContent)
	}

	tc.protocols[p] = content
}

func (tc TrafficContent) Protocol(p Protocol) *ProtocolContent {
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

	return content
}

func (tc TrafficContent) TCP() *ProtocolContent {
	return tc.Protocol(ProtocolTCP)
}

func (tc TrafficContent) UDP() *ProtocolContent {
	return tc.Protocol(ProtocolUDP)
}

func (tc TrafficContent) ICMPv4() *ProtocolContent {
	return tc.Protocol(ProtocolICMPv4)
}

func (tc TrafficContent) ICMPv6() *ProtocolContent {
	return tc.Protocol(ProtocolICMPv6)
}
