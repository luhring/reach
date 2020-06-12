package aws

import (
	"fmt"
	"net"

	"github.com/luhring/reach/reach"
	"github.com/luhring/reach/reach/traffic"
)

// FactorKindSecurityGroupRules specifies the unique name for the security group rules kind of factor.
const FactorKindSecurityGroupRules = "SecurityGroupRules"

type securityGroupRulesFactor struct {
	RuleComponents []securityGroupRulesFactorComponent
}

func (eni ElasticNetworkInterface) securityGroupRulesFactorForward(
	client DomainClient,
	leftEdge reach.Edge,
) (*reach.Factor, error) {
	sgs, err := eni.securityGroups(client)
	if err != nil {
		return nil, err
	}

	var ip net.IP
	var rules func(sg SecurityGroup) []SecurityGroupRule
	var direction securityGroupRuleDirection

	switch flow := eni.flow(leftEdge.Tuple, leftEdge.ConnectsInterface); flow {
	case reach.FlowOutbound:
		ip = leftEdge.Tuple.Dst
		rules = func(sg SecurityGroup) []SecurityGroupRule { return sg.OutboundRules }
		direction = securityGroupRuleDirectionOutbound
	case reach.FlowInbound:
		ip = leftEdge.Tuple.Src
		rules = func(sg SecurityGroup) []SecurityGroupRule { return sg.InboundRules }
		direction = securityGroupRuleDirectionInbound
	default:
		return nil, fmt.Errorf("determining security group rules factors for flow '%s' is not supported", flow)
	}

	components, err := applicableSecurityGroupRules(client, sgs, ip, rules, direction)
	if err != nil {
		return nil, err
	}

	traffic, err := trafficFromSecurityGroupRulesFactorComponents(components)
	if err != nil {
		return nil, err
	}

	factor := &reach.Factor{
		Kind:     FactorKindSecurityGroupRules,
		Resource: eni.Ref(),
		Traffic:  traffic,
		Properties: securityGroupRulesFactor{
			RuleComponents: components,
		},
	}
	return factor, nil
}

func (eni ElasticNetworkInterface) securityGroupRulesFactorReturn() reach.Factor {
	return reach.Factor{
		Kind:       FactorKindSecurityGroupRules,
		Resource:   eni.Ref(),
		Traffic:    traffic.All(),
		Properties: nil,
	}
}

func applicableSecurityGroupRules(
	client DomainClient,
	sgs []SecurityGroup,
	ip net.IP,
	rules func(sg SecurityGroup) []SecurityGroupRule,
	direction securityGroupRuleDirection,
) ([]securityGroupRulesFactorComponent, error) {
	var components []securityGroupRulesFactorComponent

	for _, sg := range sgs {
		for index, rule := range rules(sg) {
			match, err := matchSecurityGroupRule(client, rule, ip)
			if err != nil {
				return nil, fmt.Errorf("unable to get applicable security group rules: %v", err)
			}

			if match != nil {
				c := securityGroupRulesFactorComponent{
					SecurityGroupID: sg.ID,
					RuleDirection:   direction,
					RuleIndex:       index,
					Match:           *match,
					Traffic:         rule.TrafficContent,
				}
				components = append(components, c)
			}
		}
	}

	return components, nil
}

func trafficFromSecurityGroupRulesFactorComponents(components []securityGroupRulesFactorComponent) (traffic.Content, error) {
	var segments []traffic.Content
	for _, component := range components {
		segments = append(segments, component.Traffic)
	}

	tc, err := traffic.Merge(segments)
	if err != nil {
		return traffic.Content{}, err
	}
	return tc, nil
}
