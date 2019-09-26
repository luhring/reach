package reach

import (
	"encoding/json"

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
	protocols[protocol] = NewProtocolContentWithPorts(protocol, &ports)

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentForTCP(ports set.PortSet) TrafficContent {
	return NewTrafficContentForPorts(ProtocolTCP, ports)
}

func NewTrafficContentForUDP(ports set.PortSet) TrafficContent {
	return NewTrafficContentForPorts(ProtocolUDP, ports)
}

func NewTrafficContentForICMP(protocol Protocol, icmp set.ICMPSet) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	protocols[protocol] = NewProtocolContentWithICMP(protocol, &icmp)

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func NewTrafficContentForICMPv4(icmp set.ICMPSet) TrafficContent {
	return NewTrafficContentForICMP(ProtocolICMPv4, icmp)
}

func NewTrafficContentForICMPv6(icmp set.ICMPSet) TrafficContent {
	return NewTrafficContentForICMP(ProtocolICMPv6, icmp)
}

func NewTrafficContentForCustomProtocol(protocol Protocol, hasContent bool) TrafficContent {
	protocols := make(map[Protocol]*ProtocolContent)
	protocols[protocol] = NewProtocolContentForCustomProtocol(protocol, hasContent)

	return TrafficContent{
		indicator: trafficContentIndicatorUnset,
		protocols: protocols,
	}
}

func (tc TrafficContent) MarshalJSON() ([]byte, error) {
	switch tc.indicator {
	case trafficContentIndicatorNone:
		return json.Marshal("[no traffic]")
	case trafficContentIndicatorAll:
		return json.Marshal("[all traffic]")
	default:
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
}

func (tc TrafficContent) All() bool {
	return tc.indicator == trafficContentIndicatorAll
}

func (tc TrafficContent) None() bool {
	return tc.indicator == trafficContentIndicatorNone
}

func (tc TrafficContent) ForProtocol(p Protocol) *ProtocolContent {
	content := tc.protocols[p]

	if content == nil {
		if p.UsesPorts() {
			return NewProtocolContentWithPortsEmpty(p)
		}

		if p.UsesICMPTypeCodes() {
			return NewProtocolContentWithICMPEmpty(p)
		}

		// custom protocol

		return NewProtocolContentForCustomProtocolEmpty(p)
	}

	return content
}

func (tc TrafficContent) TCP() *ProtocolContent {
	return tc.ForProtocol(ProtocolTCP)
}

func (tc TrafficContent) UDP() *ProtocolContent {
	return tc.ForProtocol(ProtocolUDP)
}

func (tc TrafficContent) ICMPv4() *ProtocolContent {
	return tc.ForProtocol(ProtocolICMPv4)
}

func (tc TrafficContent) ICMPv6() *ProtocolContent {
	return tc.ForProtocol(ProtocolICMPv6)
}
