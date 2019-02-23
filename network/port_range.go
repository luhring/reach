package network

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
)

const (
	minimumPort = 0     // for TCP and UDP
	maximumPort = 65535 // for TCP and UDP
)

// PortRange keeps track of ranges of ports
type PortRange struct {
	LowPort  int64
	HighPort int64
}

func NewPortRange(lowPort, highPort int64) (*PortRange, error) {
	resultPortRange := &PortRange{
		LowPort:  lowPort,
		HighPort: highPort,
	}

	if resultPortRange.isValid() {
		return resultPortRange, nil
	}

	return nil, fmt.Errorf("unable to create port range, result was invalid: %v", resultPortRange)
}

func (r *PortRange) allPorts() bool {
	return r.LowPort == minimumPort && r.HighPort == maximumPort
}

func (r *PortRange) isValid() bool {
	if false == isValidPortNumber(r.LowPort) {
		return false
	}

	if false == isValidPortNumber(r.HighPort) {
		return false
	}

	if r.LowPort > r.HighPort {
		return false
	}

	return true
}

func (r *PortRange) describesOnlyASinglePort() bool {
	return r.LowPort == r.HighPort
}

func isValidPortNumber(portNumber int64) bool {
	if portNumber < minimumPort || portNumber > maximumPort {
		return false
	}

	return true
}

func (r *PortRange) isJuxtaposedWith(other *PortRange) bool {
	if r.isValid() && other.isValid() {
		forSorting := []*PortRange{
			r,
			other,
		}
		sortPortRanges(forSorting)

		if (forSorting[0].HighPort + 1) == forSorting[1].LowPort {
			return true
		}
	}

	return false
}

func sortPortRanges(portRanges []*PortRange) {
	sort.Slice(portRanges, func(i, j int) bool {
		return portRanges[i].LowPort < portRanges[j].LowPort
	})
}

func (r *PortRange) intersectionWith(other *PortRange) *PortRange {
	if false == r.doesIntersectWith(other) {
		return nil
	}

	intersectionLowPort := getHigherOfTwoNumbers(r.LowPort, other.LowPort)
	intersectionHighPort := getLowerOfTwoNumbers(r.HighPort, other.HighPort)

	return &PortRange{
		LowPort:  intersectionLowPort,
		HighPort: intersectionHighPort,
	}
}

func (r *PortRange) doesIntersectWith(other *PortRange) bool {
	return r.HighPort >= other.LowPort && r.LowPort <= other.HighPort
}

func (r *PortRange) mergeWith(other *PortRange) (*PortRange, error) {
	if false == r.doesIntersectWith(other) && false == r.isJuxtaposedWith(other) {
		return nil, errors.New("specified port ranges cannot be merged")
	}

	mergeResultLowPort := getLowerOfTwoNumbers(r.LowPort, other.LowPort)
	mergeResultHighPort := getHigherOfTwoNumbers(r.HighPort, other.HighPort)

	if mergeResultLowPort == minimumPort && mergeResultHighPort == maximumPort {
		mergeResultLowPort = 0
		mergeResultHighPort = 0
	}

	return &PortRange{
		LowPort:  mergeResultLowPort,
		HighPort: mergeResultHighPort,
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
			if portRanges[i].doesIntersectWith(portRanges[i-1]) {
				// merge with previous
				mergeResult, err := portRanges[i-1].mergeWith(portRanges[i])
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
				portRangeFromFirstList.intersectionWith(portRangeFromSecondList)

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

func (r *PortRange) describe() string {
	if r.allPorts() {
		return "ALL ports"
	}

	if r.describesOnlyASinglePort() {
		return strconv.FormatInt(r.LowPort, 10)
	}

	return fmt.Sprintf(
		"%d-%d",
		r.LowPort,
		r.HighPort,
	)
}
