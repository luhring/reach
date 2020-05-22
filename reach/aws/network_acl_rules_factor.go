package aws

import (
	"fmt"
	"net"
	"sort"

	"github.com/luhring/reach/reach"
)

// FactorKindNetworkACLRules specifies the unique name for the networkACLRulesFactor kind of Factor.
const FactorKindNetworkACLRules = "NetworkACLRules"

type networkACLRulesFactor struct {
	RuleComponents []networkACLRulesFactorComponent
}

func (r VPCRouter) networkACLRulesFactor(
	client DomainClient,
	subnet Subnet,
	dir NetworkACLRuleDirection,
	targetIP net.IP,
) (*reach.Factor, error) {
	nacl, err := client.NetworkACL(subnet.NetworkACLID)
	if err != nil {
		return nil, err
	}

	components, err := applicableNetworkACLRules(*nacl, dir, targetIP)
	if err != nil {
		return nil, fmt.Errorf("unable to generate network ACL rules factor: %v", err)
	}

	traffic, err := trafficFromNetworkACLRulesFactorComponents(components)
	if err != nil {
		return nil, err
	}

	factor := &reach.Factor{
		Kind:       FactorKindNetworkACLRules,
		Resource:   r.Ref(),
		Traffic:    traffic,
		Properties: networkACLRulesFactor{RuleComponents: components},
	}
	return factor, nil
}

func applicableNetworkACLRules(
	nacl NetworkACL,
	dir NetworkACLRuleDirection,
	ip net.IP,
) ([]networkACLRulesFactorComponent, error) {
	var rules []NetworkACLRule
	switch dir {
	case NetworkACLRuleDirectionOutbound:
		rules = nacl.OutboundRules
	case NetworkACLRuleDirectionInbound:
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
			if rule.Allows() {
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
