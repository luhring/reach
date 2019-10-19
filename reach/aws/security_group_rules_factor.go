package aws

import (
	"github.com/luhring/reach/reach"
)

// FactorKindSecurityGroupRules specifies the unique name for the security group rules kind of factor.
const FactorKindSecurityGroupRules = "SecurityGroupRules"

type securityGroupRuleMatchBasis string

const securityGroupRuleMatchBasisIP securityGroupRuleMatchBasis = "IP"
const securityGroupRuleMatchBasisSGRef securityGroupRuleMatchBasis = "SecurityGroupReference"

type securityGroupRulesFactor struct {
	ComponentRules []securityGroupRulesFactorComponent
}

type securityGroupRulesFactorComponent struct {
	SecurityGroup reach.ResourceReference
	RuleDirection securityGroupRuleDirection
	RuleIndex     int
	Match         securityGroupRuleMatch
}

type securityGroupRuleMatch struct {
	Basis securityGroupRuleMatchBasis
	Value interface{}
}

// String returns the string representation of a security group rule match.
func (basis securityGroupRuleMatchBasis) String() string {
	switch basis {
	case securityGroupRuleMatchBasisIP:
		return "IP address"
	case securityGroupRuleMatchBasisSGRef:
		return "attached security group"
	default:
		return "[unknown match basis]"
	}
}

func (eni ElasticNetworkInterface) newSecurityGroupRulesFactor(
	rc *reach.ResourceCollection,
	p reach.Perspective,
	awsP perspective,
	targetENI *ElasticNetworkInterface,
) (*reach.Factor, error) {
	var componentRules []securityGroupRulesFactorComponent
	var trafficContentSegments []reach.TrafficContent

	for _, id := range eni.SecurityGroupIDs {
		ref := reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindSecurityGroup,
			ID:     id,
		}

		sg := rc.Get(ref).Properties.(SecurityGroup)

		for ruleIndex, rule := range awsP.getSecurityGroupRules(sg) {
			var match *securityGroupRuleMatch

			// check ip match
			match = rule.matchByIP(p.Other.IPAddress)

			// check SG ref match (only if we don't already have a match)
			if match == nil {
				match = rule.matchBySecurityGroup(targetENI)
			}

			if match != nil {
				component := securityGroupRulesFactorComponent{
					SecurityGroup: ref,
					RuleDirection: awsP.ruleDirection,
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

	props := securityGroupRulesFactor{
		ComponentRules: componentRules,
	}

	return &reach.Factor{
		Kind:       FactorKindSecurityGroupRules,
		Resource:   eni.ToResourceReference(),
		Traffic:    tc,
		Properties: props,
	}, nil
}
