package reach

import (
	"fmt"
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

func (v *InterfaceVector) getAllowedTrafficViaSecurityGroups() []*network.TrafficAllowance {
	var outboundAllowedTraffic []*network.TrafficAllowance

	for _, securityGroup := range v.Source.SecurityGroups {
		for _, rule := range securityGroup.OutboundRules {
			// if IP ranges OR security group include DESTINATION instance, add the port range to outbound ports
			var target MatchedTarget

			doesApply, target := rule.doesApplyToInterface(v.Destination)
			if doesApply {
				fmt.Printf(
					"- source interface has a security group (%v) that has an outbound rule that allows access to destination interface's %v for the following traffic: %v\n",
					securityGroup.Name,
					target.Describe(),
					rule.TrafficAllowance.Describe(),
				)
				outboundAllowedTraffic = append(outboundAllowedTraffic, rule.TrafficAllowance)
			}
		}
	}

	outboundAllowedTraffic = network.ConsolidateTrafficAllowances(outboundAllowedTraffic)

	var inboundAllowedTraffic []*network.TrafficAllowance

	for _, securityGroup := range v.Destination.SecurityGroups {
		for _, rule := range securityGroup.InboundRules {
			// if IP ranges OR security group include SOURCE instance, add the port range to inbound ports
			var target MatchedTarget
			doesApply, target := rule.doesApplyToInterface(v.Source)
			if doesApply {
				fmt.Printf(
					"- destination interface has a security group (%v) that has an inbound rule that allows access from source interface's %v for the following traffic: %v\n",
					securityGroup.Name,
					target.Describe(),
					rule.TrafficAllowance.Describe(),
				)
				inboundAllowedTraffic = append(inboundAllowedTraffic, rule.TrafficAllowance)
			}
		}
	}

	inboundAllowedTraffic = network.ConsolidateTrafficAllowances(inboundAllowedTraffic)

	intersection := network.IntersectTrafficAllowances(outboundAllowedTraffic, inboundAllowedTraffic)

	// TODO: allow filtering by specific port (via PortRange object, probably)
	// if v.PortRange != nil {
	// 	vectorPortRangeFilter := []*network.PortRange{
	// 		v.PortRange,
	// 	}
	//
	// 	intersection = network.IntersectPortRangeSlices(intersection, vectorPortRangeFilter)
	// }

	return intersection
}
