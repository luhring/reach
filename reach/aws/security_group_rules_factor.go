package aws

import (
	"github.com/luhring/reach/reach"
)

const FactorKindSecurityGroupRules = "SecurityGroupRules"

type SecurityGroupRuleMatchBasis string

const SecurityGroupRuleMatchBasisIP SecurityGroupRuleMatchBasis = "IP"
const SecurityGroupRuleMatchBasisSGRef SecurityGroupRuleMatchBasis = "SecurityGroupReference"

type SecurityGroupRuleDirection string

const SecurityGroupRuleDirectionInbound SecurityGroupRuleDirection = "inbound"
const SecurityGroupRuleDirectionOutbound SecurityGroupRuleDirection = "outbound"

type SecurityGroupRulesFactor struct {
	ComponentRules []SecurityGroupRulesFactorComponent
}

type SecurityGroupRulesFactorComponent struct {
	SecurityGroup reach.ResourceReference
	RuleDirection SecurityGroupRuleDirection
	RuleIndex     int
	Match         SecurityGroupRuleMatch
}

type SecurityGroupRuleMatch struct {
	Basis SecurityGroupRuleMatchBasis
	Value interface{}
}

func (basis SecurityGroupRuleMatchBasis) String() string {
	switch basis {
	case SecurityGroupRuleMatchBasisIP:
		return "IP address"
	case SecurityGroupRuleMatchBasisSGRef:
		return "attached security group"
	default:
		return "[unknown match basis]"
	}
}

func (eni ElasticNetworkInterface) NewSecurityGroupRulesFactor(
	rc *reach.ResourceCollection,
	p AnalysisPerspective,
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

		for ruleIndex, rule := range p.getSecurityGroupRules(sg) {
			var match *SecurityGroupRuleMatch

			// check ip match
			match = rule.MatchByIP(p.other.IPAddress)

			// check SG ref match (only if we don't already have a match)
			if match == nil {
				match = rule.MatchBySecurityGroup(targetENI)
			}

			if match != nil {
				component := SecurityGroupRulesFactorComponent{
					SecurityGroup: ref,
					RuleDirection: p.ruleDirection,
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
		Traffic:    tc,
		Properties: props,
	}, nil
}
