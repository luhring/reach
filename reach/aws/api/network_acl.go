package api

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
	reachAWS "github.com/luhring/reach/reach/aws"
)

func (provider *ResourceProvider) GetNetworkACL(id string) (*reachAWS.NetworkACL, error) {
	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []*string{
			aws.String(id),
		},
	}
	result, err := provider.ec2.DescribeNetworkAcls(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.NetworkAcls), "network ACL", id); err != nil {
		return nil, err
	}

	networkACL := newNetworkACLFromAPI(result.NetworkAcls[0])
	return &networkACL, nil
}

func newNetworkACLFromAPI(networkACL *ec2.NetworkAcl) reachAWS.NetworkACL {
	inboundRules := getInboundNetworkACLRules(networkACL.Entries)
	outboundRules := getOutboundNetworkACLRules(networkACL.Entries)

	return reachAWS.NetworkACL{
		ID:            aws.StringValue(networkACL.NetworkAclId),
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}
}

func getNetworkACLRulesForSingleDirection(entries []*ec2.NetworkAclEntry, inbound bool) []reachAWS.NetworkACLRule {
	if entries == nil {
		return nil
	}

	rules := make([]reachAWS.NetworkACLRule, len(entries))

	for i, entry := range entries {
		if entry != nil {
			if inbound != aws.BoolValue(entry.Egress) {
				rules[i] = getNetworkACLRule(entry)
			}
		}
	}

	return rules
}

func getInboundNetworkACLRules(entries []*ec2.NetworkAclEntry) []reachAWS.NetworkACLRule {
	return getNetworkACLRulesForSingleDirection(entries, true)
}

func getOutboundNetworkACLRules(entries []*ec2.NetworkAclEntry) []reachAWS.NetworkACLRule {
	return getNetworkACLRulesForSingleDirection(entries, false)
}

func getNetworkACLRule(entry *ec2.NetworkAclEntry) reachAWS.NetworkACLRule { // note: this function ignores rule direction (inbound vs. outbound)
	if entry == nil {
		return reachAWS.NetworkACLRule{}
	}

	_, targetIPNetwork, err := net.ParseCIDR(aws.StringValue(entry.CidrBlock))
	if err != nil {
		return reachAWS.NetworkACLRule{}
	}

	var action reachAWS.NetworkACLRuleAction

	if aws.StringValue(entry.RuleAction) == ec2.RuleActionAllow {
		action = reachAWS.NetworkACLRuleActionAllow
	} else {
		action = reachAWS.NetworkACLRuleActionDeny
	}

	tc, err := newTrafficContentFromAWSNACLEntry(entry)

	if err != nil {
		panic(err) // TODO: Better error handling
	}

	return reachAWS.NetworkACLRule{
		Number:          aws.Int64Value(entry.RuleNumber),
		TrafficContent:  tc,
		TargetIPNetwork: targetIPNetwork,
		Action:          action,
	}
}

func newTrafficContentFromAWSNACLEntry(entry *ec2.NetworkAclEntry) (reach.TrafficContent, error) { // TODO: BUG! This needs to consider what rules preempt this rule, and handle set subtractions accordingly
	const errCreation = "unable to create content: %v"

	protocol, err := convertAWSIPProtocolStringToProtocol(entry.Protocol)
	if err != nil {
		return reach.TrafficContent{}, fmt.Errorf(errCreation, err)
	}

	if protocol == reach.ProtocolAll {
		return reach.NewTrafficContentForAllTraffic(), nil
	}

	if protocol.UsesPorts() {
		portSet, err := newPortSetFromAWSPortRange(entry.PortRange)
		if err != nil {
			return reach.TrafficContent{}, fmt.Errorf(errCreation, err)
		}

		return reach.NewTrafficContentForPorts(protocol, portSet), nil
	}

	if protocol == reach.ProtocolICMPv4 || protocol == reach.ProtocolICMPv6 {
		icmpSet, err := newICMPSetFromAWSICMPTypeCode(entry.IcmpTypeCode)
		if err != nil {
			return reach.TrafficContent{}, fmt.Errorf(errCreation, err)
		}

		return reach.NewTrafficContentForICMP(protocol, icmpSet), nil
	}

	return reach.NewTrafficContentForCustomProtocol(protocol, true), nil
}
