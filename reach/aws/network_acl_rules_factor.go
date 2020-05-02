package aws

import (
	"fmt"
	"net"
	"sort"

	"github.com/luhring/reach/reach"
)

const FactorKindNetworkACLRules = "NetworkACLRules"

type networkACLRulesFactor struct {
	RuleComponents []networkACLRulesFactorComponent
}

func (r VPCRouter) networkACLRulesFactor(
	client DomainClient,
	subnet Subnet,
	dir networkACLRuleDirection,
	tuple reach.IPTuple,
) (*reach.Factor, error) {
	nacl, err := client.NetworkACL(subnet.NetworkACLID)
	if err != nil {
		return nil, fmt.Errorf("unable to get factors: %v", err)
	}

	components, err := applicableNetworkACLRules(*nacl, dir, tuple.Dst)
	if err != nil {
		return nil, fmt.Errorf("unable to generate network ACL rules factor: %v", err)
	}

	traffic, err := trafficFromNetworkACLRulesFactorComponents(components)
	if err != nil {
		return nil, err
	}

	factor := &reach.Factor{
		Kind:          FactorKindNetworkACLRules,
		Resource:      r.VPC.ResourceReference(),
		Traffic:       traffic,
		ReturnTraffic: reach.TrafficContent{},
		Properties:    networkACLRulesFactor{RuleComponents: components},
	}
	return factor, nil
}

func applicableNetworkACLRules(
	nacl NetworkACL,
	dir networkACLRuleDirection,
	ip net.IP,
) ([]networkACLRulesFactorComponent, error) {
	var rules []NetworkACLRule
	switch dir {
	case networkACLRuleDirectionOutbound:
		rules = nacl.OutboundRules
	case networkACLRuleDirectionInbound:
		rules = nacl.InboundRules
	default:
		return nil, fmt.Errorf("unexpected network ACL rule direction: %s", dir)
	}

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Number < rules[j].Number
	})

	var components []networkACLRulesFactorComponent
	decided := reach.NewTrafficContentForNoTraffic()

	for _, rule := range rules {
		match := matchNetworkACLRule(rule, ip)
		if match != nil {
			effectiveTraffic, err := rule.TrafficContent.Subtract(decided)
			if err != nil {
				return nil, err
			}

			c := networkACLRulesFactorComponent{
				NetworkACLID:  nacl.ID,
				RuleDirection: dir,
				RuleIndex:     rule.Number,
				Match:         *match,
				Traffic:       effectiveTraffic,
			}
			components = append(components, c)
		}

		var err error
		decided, err = reach.NewTrafficContentFromMergingMultiple([]reach.TrafficContent{
			decided,
			rule.TrafficContent,
		})
		if err != nil {
			return nil, err
		}
	}

	return components, nil
}

func trafficFromNetworkACLRulesFactorComponents(
	components []networkACLRulesFactorComponent,
) (reach.TrafficContent, error) {
	var segments []reach.TrafficContent
	for _, c := range components {
		segments = append(segments, c.Traffic)
	}

	tc, err := reach.NewTrafficContentFromMergingMultiple(segments)
	if err != nil {
		return reach.TrafficContent{}, err

	}
	return tc, nil
}
