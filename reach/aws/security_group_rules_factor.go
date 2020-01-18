package aws

import (
	"github.com/luhring/reach/reach"
)

// FactorKindSecurityGroupRules specifies the unique name for the security group rules kind of factor.
const FactorKindSecurityGroupRules = "SecurityGroupRules"

type securityGroupRulesFactor struct {
	RuleComponents []securityGroupRulesFactorComponent
}

func (eni ElasticNetworkInterface) newSecurityGroupRulesFactor(
	rc *reach.ResourceCollection,
	p reach.Perspective,
	awsP perspective,
	targetENI ElasticNetworkInterface,
) (*reach.Factor, error) {
	var ruleComponents []securityGroupRulesFactorComponent
	var trafficContentSegments []reach.TrafficContent

	for _, id := range eni.SecurityGroupIDs {
		ref := reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindSecurityGroup,
			ID:     id,
		}

		sg := rc.Get(ref).Properties.(SecurityGroup)

		for ruleIndex, rule := range awsP.securityGroupRules(sg) {
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
					RuleDirection: awsP.securityGroupRuleDirection,
					RuleIndex:     ruleIndex,
					Match:         *match,
					Traffic:       rule.TrafficContent,
				}

				trafficContentSegments = append(trafficContentSegments, rule.TrafficContent)
				ruleComponents = append(ruleComponents, component)
			}
		}
	}

	tc, err := reach.NewTrafficContentFromMergingMultiple(trafficContentSegments)
	if err != nil {
		return nil, err
	}

	props := securityGroupRulesFactor{
		RuleComponents: ruleComponents,
	}

	return &reach.Factor{
		Kind:          FactorKindSecurityGroupRules,
		Resource:      eni.ToResourceReference(),
		Traffic:       tc,
		ReturnTraffic: reach.NewTrafficContentForAllTraffic(),
		Properties:    props,
	}, nil
}
