package set

import (
	"fmt"
)

const (
	allICMPTypes = -1
	allICMPCodes = -1

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

type IcmpSet struct {
	set set
	Version uint
}

func NewEmptyICMPSet(version uint) IcmpSet {
	validateICMPVersion(version)

	return IcmpSet{
		set: newEmptySet(),
		Version: version,
	}
}

func NewFullICMPSet(version uint) *IcmpSet {
	validateICMPVersion(version)

	return &IcmpSet{
		set: newCompleteSet(),
		Version: version,
	}
}

func NewICMPSetFromICMPType(version uint, icmpType int64) (*IcmpSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return nil, fmt.Errorf("unable to use icmpType: %v", err)
	}

	startingICMPTypeCodeIndex := encodeICMPTypeCode(uint(icmpType), minimumICMPCode)
	endingICMPTypeCodeIndex := encodeICMPTypeCode(uint(icmpType), maximumICMPCode)

	set := newSetFromRange(startingICMPTypeCodeIndex, endingICMPTypeCodeIndex)

	return &IcmpSet{
		set: set,
		Version: version,
	}, nil
}

func NewICMPSetFromICMPTypeCode(version uint, icmpType, icmpCode int64) (*IcmpSet, error) {
	if err := validateICMPType(icmpType); err != nil {
		return nil, fmt.Errorf("unable to use icmpType: %v", err)
	}

	if err := validateICMPCode(icmpCode); err != nil {
		return nil, fmt.Errorf("unable to use icmpCode: %v", err)
	}

	typeCodeIndex := encodeICMPTypeCode(uint(icmpType), uint(icmpCode))
	set := NewSetFromSingleValue(typeCodeIndex)

	return &IcmpSet{
		set: set,
		Version: version,
	}, nil
}

func (s IcmpSet) Intersect(other IcmpSet) IcmpSet {
	return IcmpSet{
		set: s.set.intersect(other.set),
	}
}

func (s IcmpSet) Merge(other IcmpSet) IcmpSet {
	return IcmpSet{
		set: s.set.merge(other.set),
	}
}

func (s IcmpSet) Subtract(other IcmpSet) IcmpSet {
	return IcmpSet{
		set: s.set.subtract(other.set),
	}
}

func (s IcmpSet) String() string {
	protocol := fmt.Sprintf("ICMPv%d", s.Version)

	if s.set.complete {
		return fmt.Sprintf("all %s types and codes", protocol)
	}

	result := ""

	// {type : [allowed codes, ...]}
	typeCodes := map[uint][]uint{}
	for value := range s.set.Iterate() {
		// todo: finalize types to prevent cast
		icmpType, icmpCode := decodeICMPTypeCode(uint16(value))
		if _, ok := typeCodes[icmpType]; !ok {
			typeCodes[icmpType] = []uint{}
		}
		typeCodes[icmpType] = append(typeCodes[icmpType], icmpCode)
	}

	idx := 0
	for icmpType, icmpCodes := range typeCodes {
		var codeStr string
		if len(icmpCodes) == maximumICMPCode+1 {
			codeStr = "all codes"
		} else if len(icmpCodes) == 1 {
			codeStr = fmt.Sprintf("code %+v", icmpCodes[0])
		} else {
			codeStr = fmt.Sprintf("codes: %+v", icmpCodes)
		}

		result += fmt.Sprintf("type %v (%v), %+v", icmpType, getTypeName(s.Version, icmpType), codeStr)
		idx++
		if idx < len(typeCodes) {
			result += "; "
		}
	}

	return result
}

func validateICMPVersion(icmpVersion uint) {
	if icmpVersion == 4 || icmpVersion == 6 {
		return
	}

	panic(fmt.Errorf("bad icmp version given: %d", icmpVersion))
}

func validateICMPType(icmpType int64) error {
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

func validateICMPCode(icmpCode int64) error {
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

	return uint16((icmpType<<bitSize) | icmpCode)
}

func decodeICMPTypeCode(value uint16) (icmpType, icmpCode uint) {
	const bitSize = 8
	icmpType = uint(value >> bitSize)
	icmpCode = uint(value & 0xff)
	return
}

func getTypeName(version, icmpType uint) string {
	if typeName := icmpv4TypeNames[icmpType]; version == 4 && typeName != "" {
		return typeName
	} else if typeName := icmpv6TypeNames[icmpType]; version == 6 && typeName != "" {
		return typeName
	}

	return "unknown type"
}




