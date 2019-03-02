package reach

import (
	"github.com/logrusorgru/aurora"
	"github.com/luhring/reach/network"
	"github.com/mgutz/ansi"
)

const (
	sourcePerspective      = 100
	destinationPerspective = 200
)

type InterfaceVector struct {
	Source      *NetworkInterface
	Destination *NetworkInterface
	PortRange   *network.PortRange
}

func (v *InterfaceVector) sameSubnet() bool {
	return v.Source.SubnetID == v.Destination.SubnetID
}

func (v *InterfaceVector) explainSourceAndDestination() Explanation {
	var explanation Explanation

	explanation.AddLineFormat("source network interface: %v", aurora.Bold(v.Source.Name))
	explanation.AddLineFormat("destination network interface: %v", aurora.Bold(v.Destination.Name))

	return explanation
}

func (v *InterfaceVector) analyzeSecurityGroups() ([]*network.TrafficAllowance, Explanation) {
	var explanation Explanation
	explanation.AddLineFormat("%v analysis", aurora.Bold("security group"))

	outboundAllowedTraffic, sourceExplanation := v.analyzeSinglePerspectiveViaSecurityGroups(sourcePerspective)
	explanation.Subsume(sourceExplanation)

	// ----

	inboundAllowedTraffic, destinationExplanation := v.analyzeSinglePerspectiveViaSecurityGroups(destinationPerspective)
	explanation.Subsume(destinationExplanation)

	intersection := network.IntersectTrafficAllowances(outboundAllowedTraffic, inboundAllowedTraffic)

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

func (v *InterfaceVector) analyzeSinglePerspectiveViaSecurityGroups(perspective int) ([]*network.TrafficAllowance, Explanation) {
	var securityGroupsExplanation Explanation

	var perspectiveDescriptor string
	var perspectiveInterface *NetworkInterface
	var observedInterface *NetworkInterface
	var observedDescriptor string
	var rulePerspective string
	var getRulesForPerspective func(sg *SecurityGroup) []*SecurityGroupRule
	if perspective == sourcePerspective {
		perspectiveDescriptor = "source"
		perspectiveInterface = v.Source
		observedInterface = v.Destination
		observedDescriptor = "destination"
		rulePerspective = "outbound"
		getRulesForPerspective = func(sg *SecurityGroup) []*SecurityGroupRule { return sg.OutboundRules }
	} else {
		perspectiveDescriptor = "destination"
		perspectiveInterface = v.Destination
		observedInterface = v.Source
		observedDescriptor = "source"
		rulePerspective = "inbound"
		getRulesForPerspective = func(sg *SecurityGroup) []*SecurityGroupRule { return sg.InboundRules }
	}

	securityGroupsExplanation.AddLineFormat("%s network interface's security groups:", perspectiveDescriptor)

	var allowedTraffic []*network.TrafficAllowance

	for _, securityGroup := range perspectiveInterface.SecurityGroups {
		var securityGroupExplanation Explanation

		securityGroupExplanation.AddLineFormat(
			"%v",
			aurora.Bold(securityGroup.LongName()),
		)
		securityGroupExplanation.AddLineFormat(
			"%s security group rules that refer to the %s network interface:",
			rulePerspective,
			observedDescriptor,
		)

		var ruleMatches []RuleMatch

		for _, rule := range getRulesForPerspective(securityGroup) {
			ruleMatch := rule.matchWithInterface(observedInterface)
			if ruleMatch != nil {
				ruleMatches = append(ruleMatches, ruleMatch)
			}
		}

		if len(ruleMatches) == 0 {
			var noMatchingRules Explanation
			noMatchingRules.AddLine(ansi.Color("(no rules that meet this criteria)", "red"))
			securityGroupExplanation.Subsume(noMatchingRules)
		}

		for _, match := range ruleMatches {
			securityGroupExplanation.Subsume(match.Explain(observedDescriptor))
			allowedTraffic = append(allowedTraffic, match.GetRule().TrafficAllowance)
		}

		securityGroupsExplanation.Subsume(securityGroupExplanation)
	}

	allowedTraffic = network.ConsolidateTrafficAllowances(allowedTraffic)

	return allowedTraffic, securityGroupsExplanation
}
