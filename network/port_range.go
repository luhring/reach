package network

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// PortRange keeps track of ranges of ports
type PortRange struct {
	DoesSpecifyAllPorts     bool
	LowPort                 int64
	HighPort                int64
	DoesSpecifyAllProtocols bool
	Protocol                string
}

const minimumPort = 1
const maximumPort = 65535

func (portRange *PortRange) isValid() bool {
	if portRange.DoesSpecifyAllPorts {
		return true
	}

	if false == isValidPortNumber(portRange.LowPort) {
		return false
	}

	if false == isValidPortNumber(portRange.HighPort) {
		return false
	}

	if portRange.LowPort > portRange.HighPort {
		return false
	}

	return true
}

func isValidPortNumber(portNumber int64) bool {
	const minimumPortNumber = 1
	const maximumPortNumber = 65535

	if portNumber < minimumPortNumber || portNumber > maximumPortNumber {
		return false
	}

	return true
}

func (portRange *PortRange) doesDescribeOnlyASinglePort() bool {
	return false == portRange.DoesSpecifyAllPorts && portRange.LowPort == portRange.HighPort
}

func arePortRangesJuxtaposed(portRanges [2]*PortRange) bool {
	if portRanges[0].DoesSpecifyAllPorts || portRanges[1].DoesSpecifyAllPorts {
		return false
	}

	if (portRanges[0].DoesSpecifyAllProtocols || portRanges[1].DoesSpecifyAllProtocols) && false == (portRanges[0].DoesSpecifyAllProtocols && portRanges[1].DoesSpecifyAllProtocols) {
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

	if portRange.DoesSpecifyAllPorts {
		return &PortRange{
			DoesSpecifyAllPorts:     false,
			LowPort:                 secondPortRange.LowPort,
			HighPort:                secondPortRange.HighPort,
			DoesSpecifyAllProtocols: portRange.DoesSpecifyAllProtocols,
			Protocol:                portRange.Protocol,
		}
	}

	if secondPortRange.DoesSpecifyAllPorts {
		return &PortRange{
			DoesSpecifyAllPorts:     false,
			LowPort:                 portRange.LowPort,
			HighPort:                portRange.HighPort,
			DoesSpecifyAllProtocols: portRange.DoesSpecifyAllProtocols,
			Protocol:                portRange.Protocol,
		}
	}

	intersectionLowPort := getHigherOfTwoNumbers(portRange.LowPort, secondPortRange.LowPort)
	intersectionHighPort := getLowerOfTwoNumbers(portRange.HighPort, secondPortRange.HighPort)

	return &PortRange{
		DoesSpecifyAllPorts:     false,
		LowPort:                 intersectionLowPort,
		HighPort:                intersectionHighPort,
		DoesSpecifyAllProtocols: portRange.DoesSpecifyAllProtocols,
		Protocol:                portRange.Protocol,
	}
}

func doPortRangesIntersect(firstPortRange *PortRange, secondPortRange *PortRange) bool {
	if firstPortRange.DoesSpecifyAllProtocols != secondPortRange.DoesSpecifyAllProtocols {
		return false
	}

	if firstPortRange.Protocol != secondPortRange.Protocol {
		return false
	}

	if firstPortRange.DoesSpecifyAllPorts || secondPortRange.DoesSpecifyAllPorts {
		return true
	}

	return firstPortRange.HighPort >= secondPortRange.LowPort && firstPortRange.LowPort <= secondPortRange.HighPort
}

func mergePortRanges(firstPortRange *PortRange, secondPortRange *PortRange) *PortRange {
	portRanges := [2]*PortRange{
		firstPortRange,
		secondPortRange,
	}

	if false == doPortRangesIntersect(firstPortRange, secondPortRange) && false == arePortRangesJuxtaposed(portRanges) {
		return nil // should replace this with more idiomatic error presentation
	}

	if firstPortRange.DoesSpecifyAllPorts || secondPortRange.DoesSpecifyAllPorts {
		return &PortRange{
			DoesSpecifyAllPorts:     true,
			LowPort:                 0,
			HighPort:                0,
			DoesSpecifyAllProtocols: firstPortRange.DoesSpecifyAllProtocols,
			Protocol:                firstPortRange.Protocol,
		}
	}

	mergeResultLowPort := getLowerOfTwoNumbers(firstPortRange.LowPort, secondPortRange.LowPort)
	mergeResultHighPort := getHigherOfTwoNumbers(firstPortRange.HighPort, secondPortRange.HighPort)
	mergeResultDoesSpecifyAllPorts := false

	if mergeResultLowPort == minimumPort && mergeResultHighPort == maximumPort {
		mergeResultLowPort = 0
		mergeResultHighPort = 0
		mergeResultDoesSpecifyAllPorts = true
	}

	return &PortRange{
		DoesSpecifyAllPorts:     mergeResultDoesSpecifyAllPorts,
		LowPort:                 mergeResultLowPort,
		HighPort:                mergeResultHighPort,
		DoesSpecifyAllProtocols: firstPortRange.DoesSpecifyAllProtocols,
		Protocol:                firstPortRange.Protocol,
	}
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
				portRanges[i-1] = mergePortRanges(portRanges[i-1], portRanges[i])
				portRanges = append(portRanges[:i], portRanges[i+1:]...)
				// start from the top
				i = 0
			}
		}
	}

	return portRanges
}

// GetIntersectionBetweenTwoListsOfPortRanges ...
func GetIntersectionBetweenTwoListsOfPortRanges(
	firstListOfPortRanges []*PortRange,
	secondListOfPortRanges []*PortRange,
) []*PortRange {
	var intersectionPortRanges []*PortRange

	for _, portRangeFromFirstList := range firstListOfPortRanges {
		for _, portRangeFromSecondList := range secondListOfPortRanges {
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
	if portRange.DoesSpecifyAllProtocols {
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
	if portRange.DoesSpecifyAllPorts {
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
