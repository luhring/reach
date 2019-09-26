package set

import (
	"fmt"
)

const (
	AllICMPTypes = -1
	AllICMPCodes = -1

	minimumICMPType = 0
	maximumICMPType = 255
	minimumICMPCode = 0
	maximumICMPCode = 255
)

var icmpv4TypeNames = map[uint]string{
	0:  "echo reply",
	1:  "reserved",
	2:  "reserved",
	3:  "destination unreachable",
	4:  "source quench",
	5:  "redirect message",
	6:  "alternate host address",
	7:  "reserved",
	8:  "echo request",
	9:  "router advertisement",
	10: "router solicitation",
	11: "time exceeded",
	12: "parameter problem: bad IP header",
	13: "timestamp",
	14: "timestamp reply",
	15: "information request",
	16: "information reply",
	17: "address mask request",
	18: "address mask reply",
	30: "information request",
	31: "datagram conversion error",
	32: "mobile host redirect",
	33: "where are you",
	34: "here I am",
	35: "mobile registration request",
	36: "mobile registration reply",
	37: "domain name request",
	38: "domain name reply",
	39: "SKIP algorithm discovery protocol",
	40: "Photuris, security failures",
}

var icmpv6TypeNames = map[uint]string{
	1:   "destination unreachable",
	2:   "packet too big",
	3:   "time exceeded",
	4:   "parameter problem",
	100: "private experimentation",
	101: "private experimentation",
	128: "echo request",
	129: "echo reply",
	130: "multicast listener query (MLD)",
	131: "multicast listener report (MLD)",
	132: "multicast listener done (MLD)",
	133: "router solicitation (NDP)",
	134: "router advertisement (NDP)",
	135: "neighbor solicitation (NDP)",
	136: "neighbor advertisement (NDP)",
	137: "redirect message",
	138: "router renumbering",
	139: "ICMP node information query",
	140: "ICMP node information response",
}

type ICMPSet struct {
	set Set
}

func NewEmptyICMPSet() ICMPSet {
	return ICMPSet{
		set: newEmptySet(),
	}
}

func NewFullICMPSet() *ICMPSet {
	return &ICMPSet{
		set: newCompleteSet(),
	}
}

func NewICMPSetFromICMPType(icmpType uint8) (*ICMPSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return nil, fmt.Errorf("unable to use icmpType: %v", err)
	}

	startingICMPTypeCodeIndex := encodeICMPTypeCode(uint(icmpType), minimumICMPCode)
	endingICMPTypeCodeIndex := encodeICMPTypeCode(uint(icmpType), maximumICMPCode)

	set := newSetFromRange(startingICMPTypeCodeIndex, endingICMPTypeCodeIndex)

	return &ICMPSet{
		set: set,
	}, nil
}

func NewICMPSetFromICMPTypeCode(icmpType, icmpCode uint8) (*ICMPSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return nil, fmt.Errorf("unable to use icmpType: %v", err)
	}

	if err := validateICMPCode(icmpCode); err != nil {
		return nil, fmt.Errorf("unable to use icmpCode: %v", err)
	}

	typeCodeIndex := encodeICMPTypeCode(uint(icmpType), uint(icmpCode))
	set := NewSetForSingleValue(typeCodeIndex)

	return &ICMPSet{
		set: set,
	}, nil
}

func (s ICMPSet) Intersect(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.intersect(other.set),
	}
}

func (s ICMPSet) Merge(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.merge(other.set),
	}
}

func (s ICMPSet) Subtract(other ICMPSet) ICMPSet {
	return ICMPSet{
		set: s.set.subtract(other.set),
	}
}

func validateICMPType(icmpType uint8) error {
	if icmpType < minimumICMPType || icmpType > maximumICMPType {
		return fmt.Errorf(
			"icmpType value %v is not valid, must be between %v and %v (inclusive)",
			icmpType,
			minimumICMPType,
			maximumICMPType,
		)
	}

	return nil
}

func validateICMPCode(icmpCode uint8) error {
	if icmpCode < minimumICMPCode || icmpCode > maximumICMPCode {
		return fmt.Errorf(
			"icmpCode value %v is not valid, must be between %v and %v (inclusive)",
			icmpCode,
			minimumICMPCode,
			maximumICMPCode,
		)
	}

	return nil
}

func encodeICMPTypeCode(icmpType, icmpCode uint) uint16 {
	const bitSize = 8

	return uint16((icmpType << bitSize) | icmpCode)
}
