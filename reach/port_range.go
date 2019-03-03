package reach

import (
	"errors"
	"fmt"
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

func newPortRange(lowPort, highPort int64) (*PortRange, error) {
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

func (r *PortRange) describe() string {
	if r.allPorts() {
		return "ALL ports"
	}

	if r.describesOnlyASinglePort() {
		return strconv.FormatInt(r.LowPort, 10)
	}

	return fmt.Sprintf(
		"%d - %d",
		r.LowPort,
		r.HighPort,
	)
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

func isValidPortNumber(portNumber int64) bool {
	if portNumber < minimumPort || portNumber > maximumPort {
		return false
	}

	return true
}
