package set

import "fmt"

var icmpv4TypeNames = map[uint8]string{
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

var icmpv6TypeNames = map[uint8]string{
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

type ICMPTypeCode struct {
	icmpType uint8
	icmpCode uint8
}

func (i ICMPTypeCode) StringV4() string {
	typeName := GetICMPv4TypeName(i.icmpType)
	return fmt.Sprintf("%s (code %d)", typeName, i.icmpCode)
}

func (i ICMPTypeCode) StringV6() string {
	typeName := GetICMPv6TypeName(i.icmpType)
	return fmt.Sprintf("%s (code %d)", typeName, i.icmpCode)
}

func GetICMPv4TypeName(icmpType uint8) string {
	typeName, exists := icmpv4TypeNames[icmpType]
	if !exists {
		typeName = fmt.Sprintf("(unnamed ICMPv4 type: %d)", icmpType)
	}

	return typeName
}

func GetICMPv6TypeName(icmpType uint8) string {
	typeName, exists := icmpv6TypeNames[icmpType]
	if !exists {
		typeName = fmt.Sprintf("(unnamed ICMPv6 type: %d)", icmpType)
	}

	return typeName
}