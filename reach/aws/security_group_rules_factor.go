package aws

import (
	"net"

	"github.com/luhring/reach/reach"
)

const FactorKindSecurityGroupRules = "SecurityGroupRules"

type SecurityGroupRuleMatchBasis string

const SecurityGroupRuleMatchBasisIP SecurityGroupRuleMatchBasis = "IP"
const SecurityGroupRuleMatchBasisSGRef SecurityGroupRuleMatchBasis = "SecurityGroupReference"

type SecurityGroupRulesFactor struct {
	ComponentRules []SecurityGroupRulesFactorComponent
}

type SecurityGroupRulesFactorComponent struct {
	SecurityGroup reach.ResourceReference
	RuleIndex     int
	Match         SecurityGroupRuleMatch
}

type SecurityGroupRuleMatch struct {
	Basis SecurityGroupRuleMatchBasis
	Value interface{}
}

func (eni ElasticNetworkInterface) NewSecurityGroupRulesFactor(
	rc *reach.ResourceCollection,
	getRules func(sg SecurityGroup) []SecurityGroupRule,
	targetIPAddress net.IP,
	targetENI *ElasticNetworkInterface,
) (*reach.Factor, error) {
	var componentRules []SecurityGroupRulesFactorComponent
	var trafficContentSegments []reach.TrafficContent

	for _, id := range eni.SecurityGroupIDs {
		ref := reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindSecurityGroup,
			ID:     id,
		}

		sg := rc.Get(ref).Properties.(SecurityGroup)

		for ruleIndex, rule := range getRules(sg) {
			var match *SecurityGroupRuleMatch

			// check ip match
			match = rule.MatchByIP(targetIPAddress)

			// check SG ref match (only if we don't already have a match)
			if match == nil {
				match = rule.MatchBySecurityGroup(targetENI)
			}

			if match != nil {
				component := SecurityGroupRulesFactorComponent{
					SecurityGroup: ref,
					RuleIndex:     ruleIndex,
					Match:         *match,
				}

				trafficContentSegments = append(trafficContentSegments, rule.TrafficContent)
				componentRules = append(componentRules, component)
			}
		}
	}

	tc, err := reach.NewTrafficContentFromMergingMultiple(trafficContentSegments)
	if err != nil {
		return nil, err
	}

	props := SecurityGroupRulesFactor{
		ComponentRules: componentRules,
	}

	return &reach.Factor{
		Kind:       FactorKindSecurityGroupRules,
		Resource:   eni.ToResourceReference(),
		Traffic:    *tc,
		Properties: props,
	}, nil
}
