package reach

import (
	"github.com/luhring/reach/network"
)

type InterfaceVector struct {
	Source      *NetworkInterface
	Destination *NetworkInterface
	PortRange   *network.PortRange
}

func (v *InterfaceVector) sameSubnet() bool {
	return v.Source.SubnetID == v.Destination.SubnetID
}

func (v *InterfaceVector) getReachablePortsViaSecurityGroups() []*network.PortRange {
	var outboundAllowedPortRanges []*network.PortRange

	for _, securityGroup := range v.Source.SecurityGroups {
		for _, rule := range securityGroup.OutboundRules {
			// if IP ranges OR security group include DESTINATION instance, add the port range to outbound ports
			if rule.doesApplyToInterface(v.Destination) {
				outboundAllowedPortRanges = append(outboundAllowedPortRanges, rule.Ports)
			}
		}
	}

	outboundAllowedPortRanges = network.DefragmentPortRanges(outboundAllowedPortRanges)

	var inboundAllowedPortRanges []*network.PortRange

	for _, securityGroup := range v.Destination.SecurityGroups {
		for _, rule := range securityGroup.InboundRules {
			// if IP ranges OR security group include SOURCE instance, add the port range to inbound ports
			if rule.doesApplyToInterface(v.Source) {
				inboundAllowedPortRanges = append(inboundAllowedPortRanges, rule.Ports)
			}
		}
	}

	inboundAllowedPortRanges = network.DefragmentPortRanges(inboundAllowedPortRanges)

	intersection := network.IntersectPortRangeSlices(outboundAllowedPortRanges, inboundAllowedPortRanges)

	if v.PortRange != nil {
		vectorPortRangeFilter := []*network.PortRange{
			v.PortRange,
		}

		intersection = network.IntersectPortRangeSlices(intersection, vectorPortRangeFilter)
	}

	return intersection
}
