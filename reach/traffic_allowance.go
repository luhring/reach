package reach

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/mgutz/ansi"
	"sort"
	"strconv"
	"strings"
)

const (
	all                       = -1
	allName                   = "all"
	icmpName                  = "ICMP"
	tcpName                   = "TCP"
	udpName                   = "UDP"
	icmpv6Name                = "ICMPv6"
	ipProtocolNumberForICMP   = 1
	ipProtocolNumberForTCP    = 6
	ipProtocolNumberForUDP    = 17
	ipProtocolNumberForICMPv6 = 58
)

type TrafficAllowance struct {
	Protocol       int64           // -1 for all protocols
	PortRange      *PortRange      // should be nil unless protocol is TCP or UDP
	ICMPConstraint *ICMPConstraint // should be nil unless protocol is ICMP or ICMPv6
}

func newTrafficAllowanceForAllTraffic() *TrafficAllowance {
	return &TrafficAllowance{
		Protocol: all,
	}
}

func newTrafficAllowanceForTCPOrUDP(protocol int64, portRange *PortRange) *TrafficAllowance {
	return &TrafficAllowance{
		Protocol:  protocol,
		PortRange: portRange,
	}
}

func NewTrafficAllowanceForTCPPort(port uint16) *TrafficAllowance {
	return &TrafficAllowance{
		Protocol: ipProtocolNumberForTCP,
		PortRange: &PortRange{
			int64(port),
			int64(port),
		},
	}
}

func newTrafficAllowanceForICMP(protocol int64, constraint *ICMPConstraint) *TrafficAllowance {
	return &TrafficAllowance{
		Protocol:       protocol,
		ICMPConstraint: constraint,
	}
}

func newTrafficAllowanceForCustomProtocol(protocol int64) *TrafficAllowance {
	return &TrafficAllowance{
		Protocol: protocol,
	}
}

func newTrafficAllowanceFromAWS(ipProtocol *string, fromPort *int64, toPort *int64) (*TrafficAllowance, error) {
	if ipProtocol == nil {
		return nil, errors.New("cannot construct traffic allowance with nil ipProtocol")
	}

	p, err := convertAWSIPProtocolStringToProtocolNumber(ipProtocol)
	if err != nil {
		return nil, fmt.Errorf("cannot construct traffic allowance because conversion to protocol number failed: %v", err)
	}

	protocol := *p

	if protocol == all {
		return newTrafficAllowanceForAllTraffic(), nil
	}

	if protocol == ipProtocolNumberForTCP || protocol == ipProtocolNumberForUDP {
		if fromPort == nil {
			return nil, errors.New("cannot construct traffic allowance with nil fromPort")
		}

		if toPort == nil {
			return nil, errors.New("cannot construct traffic allowance with nil toPort")
		}

		portRange, err := newPortRange(*fromPort, *toPort)
		if err != nil {
			return nil, fmt.Errorf("unable to construct traffic allowance due to error constructing port range: %v", err)
		}

		return newTrafficAllowanceForTCPOrUDP(protocol, portRange), nil
	}

	if protocol == ipProtocolNumberForICMP || protocol == ipProtocolNumberForICMPv6 {
		var constraint *ICMPConstraint

		if protocol == ipProtocolNumberForICMP {
			if fromPort == nil {
				return nil, errors.New("cannot construct traffic allowance with nil fromPort")
			}

			if toPort == nil {
				return nil, errors.New("cannot construct traffic allowance with nil toPort")
			}

			constraint = &ICMPConstraint{
				Type: *fromPort,
				Code: *toPort,
			}
		} else if protocol == ipProtocolNumberForICMPv6 {
			if fromPort == nil || toPort == nil {
				constraint = &ICMPConstraint{
					Type: all,
					Code: all,
					V6:   true,
				}
			} else {
				constraint = &ICMPConstraint{
					Type: *fromPort,
					Code: *toPort,
					V6:   true,
				}
			}
		}

		return newTrafficAllowanceForICMP(protocol, constraint), nil
	}

	return newTrafficAllowanceForCustomProtocol(protocol), nil
}

func convertAWSIPProtocolStringToProtocolNumber(ipProtocol *string) (*int64, error) {
	if ipProtocol == nil {
		return nil, errors.New("unexpected nil ipProtocol")
	}

	protocolString := strings.ToLower(aws.StringValue(ipProtocol))

	if protocol, err := strconv.ParseInt(protocolString, 10, 64); err == nil {
		return &protocol, nil
	}

	var protocolNumber int64

	switch protocolString {
	case strings.ToLower(tcpName):
		protocolNumber = ipProtocolNumberForTCP
	case strings.ToLower(udpName):
		protocolNumber = ipProtocolNumberForUDP
	case strings.ToLower(icmpName):
		protocolNumber = ipProtocolNumberForICMP
	case strings.ToLower(icmpv6Name):
		protocolNumber = ipProtocolNumberForICMPv6
	default:
		return nil, errors.New("unrecognized ipProtocol value")
	}

	return aws.Int64(protocolNumber), nil
}

func (t *TrafficAllowance) allProtocols() bool {
	return t.Protocol == all
}

func (t *TrafficAllowance) allTraffic() bool {
	return t.allProtocols()
}

func (t *TrafficAllowance) usesNamedProtocol() bool {
	name := t.getProtocolName()
	return strings.EqualFold(name, tcpName) ||
		strings.EqualFold(name, udpName) ||
		strings.EqualFold(name, icmpName) ||
		strings.EqualFold(name, icmpv6Name)
}

func (t *TrafficAllowance) specifiesTCPOrUDPPortRange() bool {
	return (t.Protocol == ipProtocolNumberForTCP || t.Protocol == ipProtocolNumberForUDP) && t.PortRange != nil
}

func (t *TrafficAllowance) specifiesICMPConstraint() bool {
	return (t.Protocol == ipProtocolNumberForICMP || t.Protocol == ipProtocolNumberForICMPv6) && t.ICMPConstraint != nil
}

func (t *TrafficAllowance) getProtocolName() string {
	switch t.Protocol {
	case all:
		return allName
	case ipProtocolNumberForICMP:
		return icmpName
	case ipProtocolNumberForTCP:
		return tcpName
	case ipProtocolNumberForUDP:
		return udpName
	case ipProtocolNumberForICMPv6:
		return icmpv6Name
	default:
		return string(t.Protocol)
	}
}

func (t *TrafficAllowance) intersectWith(other *TrafficAllowance) *TrafficAllowance {
	if t.allTraffic() && other.allTraffic() {
		return newTrafficAllowanceForAllTraffic()
	}

	if t.allTraffic() {
		return other
	}

	if other.allTraffic() {
		return t
	}

	if t.Protocol != other.Protocol {
		return nil
	}

	// traffic allowances use the same protocol

	if t.specifiesTCPOrUDPPortRange() {
		portRangeIntersection := t.PortRange.intersectionWith(other.PortRange)
		if portRangeIntersection == nil {
			return nil
		}

		return newTrafficAllowanceForTCPOrUDP(t.Protocol, portRangeIntersection)
	}

	if t.specifiesICMPConstraint() {
		icmpConstraintIntersection := t.ICMPConstraint.IntersectionWith(other.ICMPConstraint)
		if icmpConstraintIntersection == nil {
			return nil
		}

		return newTrafficAllowanceForICMP(t.Protocol, icmpConstraintIntersection)
	}

	return newTrafficAllowanceForCustomProtocol(t.Protocol)
}

func (t *TrafficAllowance) mergeWith(other *TrafficAllowance) (*TrafficAllowance, error) {
	if intersection := t.intersectWith(other); intersection == nil {
		return nil, errors.New("traffic allowances cannot be merged")
	}

	if t.allTraffic() || other.allTraffic() {
		return newTrafficAllowanceForAllTraffic(), nil
	}

	// neither allow all traffic, but protocols are the same since we're past the intersection test

	if t.specifiesTCPOrUDPPortRange() {
		mergedPortRange, err := t.PortRange.mergeWith(other.PortRange)
		if err != nil {
			return nil, fmt.Errorf("unable to merge traffic allowances: %v", err)
		}

		return newTrafficAllowanceForTCPOrUDP(t.Protocol, mergedPortRange), nil
	}

	if t.specifiesICMPConstraint() {
		mergedICMPConstraint, err := t.ICMPConstraint.mergeWith(other.ICMPConstraint)
		if err != nil {
			return nil, fmt.Errorf("unable to merge traffic allowances: %v", err)
		}

		return newTrafficAllowanceForICMP(t.Protocol, mergedICMPConstraint), nil
	}

	// Not using TCP, UDP, ICMP, or ICMPv6, which means this is a custom IP protocol (shared by both allowances)
	// This also means all traffic is allowed for this custom protocol, according to AWS rules

	return newTrafficAllowanceForCustomProtocol(t.Protocol), nil
}

func (t *TrafficAllowance) describe() string {
	if t.allProtocols() {
		return "ALL traffic"
	}

	var constraintPredicate string

	if t.specifiesTCPOrUDPPortRange() {
		constraintPredicate = t.PortRange.describe()
	} else if t.specifiesICMPConstraint() {
		if t.ICMPConstraint != nil {
			constraintPredicate = fmt.Sprintf("(%v)", t.ICMPConstraint.Describe())
		}
	} else {
		constraintPredicate = "[unknown]"
	}

	return t.describeProtocol() + " " + constraintPredicate
}

func (t *TrafficAllowance) describeProtocol() string {
	if t.allProtocols() {
		return "(ALL protocols)"
	}

	if t.usesNamedProtocol() {
		return t.getProtocolName()
	}

	return fmt.Sprintf(
		"(IP protocol %v)",
		t.Protocol,
	)
}

func consolidateTrafficAllowances(allowances []*TrafficAllowance) []*TrafficAllowance {
	if allowances == nil {
		return nil
	}

	if len(allowances) == 1 {
		return allowances
	}

	sortTrafficAllowances(allowances)

	for i := 0; i < len(allowances); i++ {
		if i > 0 {
			if allowances[i].Protocol == allowances[i-1].Protocol {
				mergeResult, err := allowances[i-1].mergeWith(allowances[i])
				if err != nil {
					// we can't merge these two particular allowances... that's fine
					continue
				}

				allowances[i-1] = mergeResult
				allowances = append(allowances[:i], allowances[i+1:]...)

				// start again from beginning
				i = 0
			}
		}
	}

	return allowances
}

func intersectTrafficAllowances(
	firstTrafficAllowanceSlice []*TrafficAllowance,
	secondTrafficAllowanceSlice []*TrafficAllowance,
) []*TrafficAllowance {
	var intersectionTrafficAllowances []*TrafficAllowance

	for _, allowanceFromFirstList := range firstTrafficAllowanceSlice {
		for _, allowanceFromSecondList := range secondTrafficAllowanceSlice {
			currentIntersection := allowanceFromFirstList.intersectWith(allowanceFromSecondList)

			if currentIntersection != nil {
				intersectionTrafficAllowances = append(intersectionTrafficAllowances, currentIntersection)
			}
		}
	}

	return consolidateTrafficAllowances(intersectionTrafficAllowances)
}

func sortTrafficAllowances(allowances []*TrafficAllowance) {
	sort.Slice(allowances, func(i, j int) bool {
		if allowances[i].Protocol == allowances[j].Protocol && allowances[i].specifiesTCPOrUDPPortRange() {
			return allowances[i].PortRange.LowPort < allowances[j].PortRange.LowPort
		}

		return allowances[i].Protocol < allowances[j].Protocol
	})
}

func describeListOfTrafficAllowances(allowances []*TrafficAllowance) string {
	if len(allowances) < 1 {
		return ansi.Color("not reachable", "red") + "\n"
	}

	var description string

	for _, allowance := range allowances {
		description += "✔ " + allowance.describe() + "\n"
	}

	return description
}
