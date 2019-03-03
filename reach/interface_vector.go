package reach

import (
	"fmt"
	"github.com/mgutz/ansi"
)

type InterfaceVector struct {
	Source      *NetworkInterface
	Destination *NetworkInterface
}

func (v *InterfaceVector) sameSubnet() bool {
	return v.Source.SubnetID == v.Destination.SubnetID
}

func (v *InterfaceVector) explainSourceAndDestination() Explanation {
	explanation := newExplanation(
		fmt.Sprintf("source network interface: %v", ansi.Color(v.Source.Name, "default+b")),
		fmt.Sprintf("destination network interface: %v", ansi.Color(v.Destination.Name, "default+b")),
	)

	return explanation
}

func (v *InterfaceVector) analyzeSecurityGroups(filter *TrafficAllowance) ([]*TrafficAllowance, Explanation) {
	if filter == nil {
		filter = newTrafficAllowanceForAllTraffic()
	}

	explanation := newExplanation(
		fmt.Sprintf("%v analysis", ansi.Color("security group", "default+b")),
	)

	p := newPerspectiveFromSource(v)
	outboundAllowedTraffic, sourceExplanation := v.analyzeSinglePerspectiveViaSecurityGroups(p)
	explanation.subsume(sourceExplanation)

	// ----

	p = newPerspectiveFromDestination(v)
	inboundAllowedTraffic, destinationExplanation := v.analyzeSinglePerspectiveViaSecurityGroups(p)
	explanation.subsume(destinationExplanation)

	intersection := intersectTrafficAllowances(outboundAllowedTraffic, inboundAllowedTraffic)

	// TODO: allow filtering by specific port (via PortRange object, probably)
	// if v.PortRange != nil {
	// 	vectorPortRangeFilter := []*network.PortRange{
	// 		v.PortRange,
	// 	}
	//
	// 	intersection = network.IntersectPortRangeSlices(intersection, vectorPortRangeFilter)
	// }

	return intersection, explanation
}

func (v *InterfaceVector) analyzeSinglePerspectiveViaSecurityGroups(p perspective) ([]*TrafficAllowance, Explanation) {
	securityGroupsExplanation := newExplanation(
		fmt.Sprintf("%s network interface's security groups:", p.self),
	)

	var allowedTraffic []*TrafficAllowance

	for _, securityGroup := range p.selfInterface.SecurityGroups {
		securityGroupExplanation := newExplanation(
			ansi.Color(securityGroup.longName(), "default+b"),
		)

		securityGroupExplanation.addLineFormat(
			"%s security group rules that refer to the %s network interface:",
			p.direction,
			p.other,
		)

		var ruleMatches []RuleMatch

		for _, rule := range p.rules(securityGroup) {
			ruleMatch := rule.matchWithInterface(p.otherInterface)
			if ruleMatch != nil {
				ruleMatches = append(ruleMatches, ruleMatch)
			}
		}

		for _, match := range ruleMatches {
			securityGroupExplanation.subsume(match.explain(p.other))
			allowedTraffic = append(allowedTraffic, match.getRule().TrafficAllowance)
		}

		if len(ruleMatches) >= 1 {
			securityGroupsExplanation.subsume(securityGroupExplanation)
		}
	}

	if len(allowedTraffic) == 0 {
		noMatchingRules := newExplanation(
			fmt.Sprintf(
				ansi.Color("This network interface has no security groups with %v rules that refer to the %s network interface.", "red"),
				p.direction,
				p.other,
			),
		)

		securityGroupsExplanation.subsume(noMatchingRules)
	}

	allowedTraffic = consolidateTrafficAllowances(allowedTraffic)

	return allowedTraffic, securityGroupsExplanation
}
