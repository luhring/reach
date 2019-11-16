package aws

import (
	"fmt"
	"sort"

	"github.com/luhring/reach/reach"
)

// FactorKindNetworkACLRules specifies the unique name for the network ACL rules kind of factor.
const FactorKindNetworkACLRules = "NetworkACLRules"

const newNetworkACLRulesFactorErrFmt = "unable to compute network ACL rules factor: %v"

type networkACLRulesFactor struct {
	RuleComponentsForwardDirection []networkACLRulesFactorComponent
	RuleComponentsReturnDirection  []networkACLRulesFactorComponent
}

func (eni ElasticNetworkInterface) newNetworkACLRulesFactor(
	rc *reach.ResourceCollection,
	p reach.Perspective,
	awsP perspective,
	targetENI *ElasticNetworkInterface,
) (*reach.Factor, error) {
	subnetResource := rc.Get(reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindSubnet,
		ID:     eni.SubnetID,
	})
	if subnetResource == nil {
		return nil, fmt.Errorf("couldn't find subnet: %s", eni.SubnetID)
	}
	subnet := subnetResource.Properties.(Subnet)

	ref := reach.ResourceReference{
		Domain: ResourceDomainAWS,
		Kind:   ResourceKindNetworkACL,
		ID:     subnet.NetworkACLID,
	}

	networkACLResource := rc.Get(ref)
	if networkACLResource == nil {
		return nil, fmt.Errorf("couldn't find network ACL: %s", subnet.NetworkACLID)
	}
	networkACL := networkACLResource.Properties.(NetworkACL)

	forwardTraffic, forwardComponents, err := networkACL.effectOnForwardTraffic(p, awsP)
	if err != nil {
		return nil, fmt.Errorf(newNetworkACLRulesFactorErrFmt, err)
	}

	returnTraffic, returnComponents, err := networkACL.effectOnReturnTraffic(p, awsP)
	if err != nil {
		return nil, fmt.Errorf(newNetworkACLRulesFactorErrFmt, err)
	}

	props := networkACLRulesFactor{
		RuleComponentsForwardDirection: forwardComponents,
		RuleComponentsReturnDirection:  returnComponents,
	}

	return &reach.Factor{
		Kind:          FactorKindNetworkACLRules,
		Resource:      eni.ToResourceReference(),
		Traffic:       forwardTraffic,
		ReturnTraffic: returnTraffic,
		Properties:    props,
	}, nil
}

func (nacl NetworkACL) effectOnForwardTraffic(p reach.Perspective, awsP perspective) (reach.TrafficContent, []networkACLRulesFactorComponent, error) {
	return nacl.factorComponents(awsP.networkACLRuleDirectionForForwardTraffic, p, awsP)
}

func (nacl NetworkACL) effectOnReturnTraffic(p reach.Perspective, awsP perspective) (reach.TrafficContent, []networkACLRulesFactorComponent, error) {
	return nacl.factorComponents(awsP.networkACLRuleDirectionForReturnTraffic, p, awsP)
}

func (nacl NetworkACL) rulesForDirection(direction networkACLRuleDirection) []NetworkACLRule {
	if direction == networkACLRuleDirectionOutbound {
		return nacl.OutboundRules
	}

	return nacl.InboundRules
}

func (nacl NetworkACL) factorComponents(direction networkACLRuleDirection, p reach.Perspective, awsP perspective) (reach.TrafficContent, []networkACLRulesFactorComponent, error) {
	rules := nacl.rulesForDirection(direction)

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Number < rules[j].Number
	})

	var trafficContentSegments []reach.TrafficContent
	var ruleComponents []networkACLRulesFactorComponent
	decidedTraffic := reach.NewTrafficContentForNoTraffic()

	for _, rule := range rules {
		// Make sure rule matches
		match := rule.matchByIP(p.Other.IPAddress)
		if match == nil {
			continue // this rule doesn't match
		}

		if rule.Allows() {
			// Determine what subset of rule traffic affects outcome
			effectiveTraffic, err := rule.TrafficContent.Subtract(decidedTraffic)
			if err != nil {
				return reach.TrafficContent{}, nil, fmt.Errorf(newNetworkACLRulesFactorErrFmt, err)
			}

			// add the allowed traffic to the trafficContentSegments
			trafficContentSegments = append(trafficContentSegments, effectiveTraffic)

			// add to ruleComponents for the explanation
			ruleComponents = append(ruleComponents, networkACLRulesFactorComponent{
				NetworkACL:    nacl.ToResourceReference(),
				RuleDirection: direction,
				RuleNumber:    rule.Number,
				Match:         *match,
				Traffic:       effectiveTraffic,
			})
		}

		var err error
		decidedTraffic, err = reach.NewTrafficContentFromMergingMultiple(
			[]reach.TrafficContent{
				decidedTraffic,
				rule.TrafficContent,
			},
		)
		if err != nil {
			return reach.TrafficContent{}, nil, fmt.Errorf(newNetworkACLRulesFactorErrFmt, err)
		}
	}

	traffic, err := reach.NewTrafficContentFromMergingMultiple(trafficContentSegments)
	if err != nil {
		return reach.TrafficContent{}, nil, fmt.Errorf(newNetworkACLRulesFactorErrFmt, err)
	}

	return traffic, ruleComponents, nil
}
