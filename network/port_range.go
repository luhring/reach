package network

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

const (
	minimumPort    = 0
	maximumPort    = 65535
	tcpName        = "tcp"
	udpName        = "udp"
	icmpName       = "icmp"
	icmpv6Name     = "icmpv6"
	allIPProtocols = "-1"
)

// PortRange keeps track of ranges of ports
type PortRange struct {
	LowPort  int64
	HighPort int64
	Protocol string
}

func NewPortRange(protocol string, lowPort, highPort int64) (*PortRange, error) {
	if protocol == allIPProtocols {
		return &PortRange{
			Protocol: allIPProtocols,
			LowPort:  minimumPort,
			HighPort: maximumPort,
		}, nil
	}

	// Convert protocol to named protocol if valid protocol number is specified in string
	if i, err := strconv.ParseInt(protocol, 10, 64); err == nil {
		protocol = getIPProtocolFromNumber(i)
	}

	resultPortRange := &PortRange{
		Protocol: protocol,
		LowPort:  lowPort,
		HighPort: highPort,
	}

	if resultPortRange.isValid() {
		return resultPortRange, nil
	}

	return nil, fmt.Errorf("unable to create port range, result was invalid: %v", resultPortRange)
}

func (portRange *PortRange) IncludesAllProtocols() bool {
	return portRange.Protocol == allIPProtocols
}

func (portRange *PortRange) doesIPProtocolImplyAllPorts() bool {
	switch portRange.Protocol {
	case allIPProtocols, tcpName, udpName, icmpName, icmpv6Name:
		return false
	default:
		return true
	}
}

func (portRange *PortRange) IncludesAllPorts() bool {
	if portRange.doesIPProtocolImplyAllPorts() {
		return true
	}

	if portRange.Protocol == tcpName || portRange.Protocol == udpName {
		return portRange.LowPort == minimumPort && portRange.HighPort == maximumPort
	}

	return false
}

func getIPProtocolFromNumber(protocolNumber int64) string {
	switch protocolNumber {
	case 1:
		return icmpName
	case 6:
		return tcpName
	case 17:
		return udpName
	case 58:
		return icmpv6Name
	default:
		return string(protocolNumber)
	}
}

func (portRange *PortRange) isValid() bool {
	if portRange.IncludesAllPorts() {
		return true
	}

	if portRange.Protocol == tcpName || portRange.Protocol == udpName {
		if false == isValidPortNumber(portRange.LowPort) {
			return false
		}

		if false == isValidPortNumber(portRange.HighPort) {
			return false
		}

		if portRange.LowPort > portRange.HighPort {
			return false
		}
	}

	return true
}

func isValidPortNumber(portNumber int64) bool {
	if portNumber < minimumPort || portNumber > maximumPort {
		return false
	}

	return true
}

func (portRange *PortRange) doesDescribeOnlyASinglePort() bool {
	return false == portRange.IncludesAllPorts() && portRange.LowPort == portRange.HighPort
}

func arePortRangesJuxtaposed(portRanges [2]*PortRange) bool {
	if portRanges[0].IncludesAllPorts() || portRanges[1].IncludesAllPorts() {
		return false
	}

	if (portRanges[0].IncludesAllProtocols() || portRanges[1].IncludesAllProtocols()) && false == (portRanges[0].IncludesAllProtocols() && portRanges[1].IncludesAllProtocols()) {
		return false
	}

	if portRanges[0].Protocol != portRanges[1].Protocol {
		return false
	}

	if portRanges[0].isValid() && portRanges[1].isValid() {
		sortPortRanges(portRanges[:])

		if (portRanges[0].HighPort + 1) == portRanges[1].LowPort {
			return true
		}
	}

	return false
}

func (portRange *PortRange) getIntersection(secondPortRange *PortRange) *PortRange {
	if false == doPortRangesIntersect(portRange, secondPortRange) {
		return nil
	}

	if portRange.IncludesAllPorts() {
		return &PortRange{
			LowPort:  secondPortRange.LowPort,
			HighPort: secondPortRange.HighPort,
			Protocol: portRange.Protocol,
		}
	}

	if secondPortRange.IncludesAllPorts() {
		return &PortRange{
			LowPort:  portRange.LowPort,
			HighPort: portRange.HighPort,
			Protocol: portRange.Protocol,
		}
	}

	intersectionLowPort := getHigherOfTwoNumbers(portRange.LowPort, secondPortRange.LowPort)
	intersectionHighPort := getLowerOfTwoNumbers(portRange.HighPort, secondPortRange.HighPort)

	return &PortRange{
		LowPort:  intersectionLowPort,
		HighPort: intersectionHighPort,
		Protocol: portRange.Protocol,
	}
}

func doPortRangesIntersect(firstPortRange *PortRange, secondPortRange *PortRange) bool {
	if firstPortRange.IncludesAllProtocols() != secondPortRange.IncludesAllProtocols() {
		return false
	}

	if firstPortRange.Protocol != secondPortRange.Protocol {
		return false
	}

	if firstPortRange.IncludesAllPorts() || secondPortRange.IncludesAllPorts() {
		return true
	}

	return firstPortRange.HighPort >= secondPortRange.LowPort && firstPortRange.LowPort <= secondPortRange.HighPort
}

func mergePortRanges(firstPortRange *PortRange, secondPortRange *PortRange) (*PortRange, error) {
	portRanges := [2]*PortRange{
		firstPortRange,
		secondPortRange,
	}

	if false == doPortRangesIntersect(firstPortRange, secondPortRange) && false == arePortRangesJuxtaposed(portRanges) {
		return nil, errors.New("specified port ranges cannot be merged")
	}

	if firstPortRange.IncludesAllPorts() || secondPortRange.IncludesAllPorts() {
		return &PortRange{
			LowPort:  0,
			HighPort: 0,
			Protocol: firstPortRange.Protocol,
		}, nil
	}

	mergeResultLowPort := getLowerOfTwoNumbers(firstPortRange.LowPort, secondPortRange.LowPort)
	mergeResultHighPort := getHigherOfTwoNumbers(firstPortRange.HighPort, secondPortRange.HighPort)

	if mergeResultLowPort == minimumPort && mergeResultHighPort == maximumPort {
		mergeResultLowPort = 0
		mergeResultHighPort = 0
	}

	return &PortRange{
		LowPort:  mergeResultLowPort,
		HighPort: mergeResultHighPort,
		Protocol: firstPortRange.Protocol,
	}, nil
}

func getHigherOfTwoNumbers(firstNumber int64, secondNumber int64) int64 {
	if firstNumber > secondNumber {
		return firstNumber
	}

	return secondNumber
}

func getLowerOfTwoNumbers(firstNumber int64, secondNumber int64) int64 {
	if firstNumber < secondNumber {
		return firstNumber
	}

	return secondNumber
}

func sortPortRanges(portRanges []*PortRange) {
	sort.Slice(portRanges, func(i, j int) bool {
		return portRanges[i].LowPort < portRanges[j].LowPort
	})
}

func arePortRangesSlicesEqual(first []*PortRange, second []*PortRange) bool {
	if first == nil && second == nil {
		return true
	}

	if first == nil || second == nil {
		return false
	}

	if len(first) != len(second) {
		return false
	}

	for i := range first {
		if *first[i] != *second[i] {
			return false
		}
	}

	return true
}

// DefragmentPortRanges ...
func DefragmentPortRanges(portRanges []*PortRange) []*PortRange {
	if len(portRanges) == 1 {
		return portRanges
	}

	sortPortRanges(portRanges)

	for i := 0; i < len(portRanges); i++ {
		if i > 0 {
			if doPortRangesIntersect(portRanges[i], portRanges[i-1]) {
				// merge with previous
				mergeResult, err := mergePortRanges(portRanges[i-1], portRanges[i])
				if err != nil {
					log.Println("warning: attempted to merge unmergeable port ranges")
					continue
				}

				portRanges[i-1] = mergeResult
				portRanges = append(portRanges[:i], portRanges[i+1:]...)

				// start from the top
				i = 0
			}
		}
	}

	return portRanges
}

// IntersectPortRangeSlices ...
func IntersectPortRangeSlices(
	firstPortRangeSlice []*PortRange,
	secondPortRangeSlice []*PortRange,
) []*PortRange {
	var intersectionPortRanges []*PortRange

	for _, portRangeFromFirstList := range firstPortRangeSlice {
		for _, portRangeFromSecondList := range secondPortRangeSlice {
			currentIntersection :=
				portRangeFromFirstList.getIntersection(portRangeFromSecondList)

			if currentIntersection != nil {
				intersectionPortRanges = append(intersectionPortRanges, currentIntersection)
			}
		}
	}

	return DefragmentPortRanges(intersectionPortRanges)
}

// DescribeListOfPortRanges ...
func DescribeListOfPortRanges(listOfPortRanges []*PortRange) string {
	description := ""

	for _, portRange := range listOfPortRanges {
		description += portRange.describe() + "\n"
	}

	return description
}

func (portRange *PortRange) describe() string {
	return portRange.describeProtocol() + " " + portRange.describePorts()
}

func (portRange *PortRange) describeProtocol() string {
	if portRange.IncludesAllProtocols() {
		return "(ALL protocols)"
	}

	const textForTCP = "TCP"
	const textForUDP = "UDP"
	const textForICMP = "ICMP"
	const textForICMPv6 = "ICMPv6"
	const protocolNumberForICMPv6 = "58"

	if strings.EqualFold(portRange.Protocol, textForTCP) {
		return textForTCP
	}

	if strings.EqualFold(portRange.Protocol, textForUDP) {
		return textForUDP
	}

	if strings.EqualFold(portRange.Protocol, textForICMP) {
		return textForICMP
	}

	if strings.EqualFold(portRange.Protocol, textForICMPv6) || portRange.Protocol == protocolNumberForICMPv6 {
		return textForICMPv6
	}

	return fmt.Sprintf(
		"(IP protocol %s)",
		portRange.Protocol,
	)
}

func (portRange *PortRange) describePorts() string {
	if portRange.IncludesAllPorts() {
		return "ALL ports"
	}

	if portRange.doesDescribeOnlyASinglePort() {
		return strconv.FormatInt(portRange.LowPort, 10)
	}

	return fmt.Sprintf(
		"%d - %d",
		portRange.LowPort,
		portRange.HighPort,
	)
}
