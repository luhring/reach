package aws

import (
	"github.com/luhring/reach/reach"
)

// FactorKindSecurityGroupRules specifies the unique name for the security group rules kind of factor.
const FactorKindSecurityGroupRules = "SecurityGroupRules"

type securityGroupRulesFactor struct {
	RuleComponents []securityGroupRulesFactorComponent
}

type securityGroupRuleMatcher func(r SecurityGroupRule, other reach.NetworkPoint) *securityGroupRuleMatch

// securityGroupRulesFactorForInterDomain calculates a SecurityGroupRules factor by assuming the other network point is not within the AWS domain.
func (eni ElasticNetworkInterface) securityGroupRulesFactorForInterDomain(
	rc *reach.ResourceCollection,
	awsPerspective perspective,
	otherNetworkPoint reach.NetworkPoint,
) (*reach.Factor, error) {
	matcher := func(r SecurityGroupRule, other reach.NetworkPoint) *securityGroupRuleMatch {
		return r.matchIP(other.IPAddress)
	}

	return eni.securityGroupRulesFactor(rc, awsPerspective, otherNetworkPoint, matcher)
}

// securityGroupRulesFactorForAWSDomain calculates a SecurityGroupRules factor by assuming the other network point is within the AWS domain.
func (eni ElasticNetworkInterface) securityGroupRulesFactorForAWSDomain(
	rc *reach.ResourceCollection,
	awsPerspective perspective,
	otherNetworkPoint reach.NetworkPoint,
) (*reach.Factor, error) {
	matcher := func(r SecurityGroupRule, other reach.NetworkPoint) *securityGroupRuleMatch {
		match := r.matchIP(other.IPAddress)
		if match != nil {
			return match
		}

		if eni := ElasticNetworkInterfaceFromNetworkPoint(other, rc); eni != nil {
			return r.matchSecurityGroupAttachedToENI(*eni)
		}

		return nil
	}

	return eni.securityGroupRulesFactor(rc, awsPerspective, otherNetworkPoint, matcher)
}

func (eni ElasticNetworkInterface) securityGroupRulesFactor(
	rc *reach.ResourceCollection,
	awsPerspective perspective,
	otherNetworkPoint reach.NetworkPoint,
	matcher securityGroupRuleMatcher,
) (*reach.Factor, error) {
	var components []securityGroupRulesFactorComponent
	var trafficContentSegments []reach.TrafficContent

	for _, sgID := range eni.SecurityGroupIDs {
		sg := rc.Get(reach.ResourceReference{
			Domain: ResourceDomainAWS,
			Kind:   ResourceKindSecurityGroup,
			ID:     sgID,
		}).Properties.(SecurityGroup)

		for i, rule := range awsPerspective.securityGroupRules(sg) {
			if match := matcher(rule, otherNetworkPoint); match != nil {
				trafficContentSegments = append(trafficContentSegments, rule.TrafficContent)

				components = append(components, securityGroupRulesFactorComponent{
					SecurityGroup: sg.ToResourceReference(),
					RuleDirection: awsPerspective.securityGroupRuleDirection,
					RuleIndex:     i,
					Match:         *match,
					Traffic:       rule.TrafficContent,
				})
			}
		}
	}

	tc, err := reach.NewTrafficContentFromMergingMultiple(trafficContentSegments)
	if err != nil {
		return nil, err
	}

	return &reach.Factor{
		Kind:          FactorKindSecurityGroupRules,
		Resource:      eni.ToResourceReference(),
		Traffic:       tc,
		ReturnTraffic: reach.NewTrafficContentForAllTraffic(),
		Properties: securityGroupRulesFactor{
			components,
		},
	}, nil
}
