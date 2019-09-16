package api

import (
	"net"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
	reachAWS "github.com/luhring/reach/reach/aws"
)

func (getter *ResourceGetter) GetNetworkACL(id string) (*reachAWS.NetworkACL, error) {
	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []*string{
			aws.String(id),
		},
	}
	result, err := getter.ec2.DescribeNetworkAcls(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(result.NetworkAcls, "network ACL", id); err != nil {
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

	trafficContent := getTrafficContentFromNetworkACLRule(entry)

	return reachAWS.NetworkACLRule{
		Number:          aws.Int64Value(entry.RuleNumber),
		TrafficContent:  trafficContent,
		TargetIPNetwork: targetIPNetwork,
		Action:          action,
	}
}

func getTrafficContentFromNetworkACLRule(entry *ec2.NetworkAclEntry) reach.TrafficContent {
	ipProtocolString := aws.StringValue(entry.Protocol)
	if ipProtocolString == "" {
		return reach.TrafficContent{}
	}

	ipProtocol, err := strconv.Atoi(ipProtocolString)
	if err != nil {
		return reach.TrafficContent{}
	}

	if ipProtocol == reach.ProtocolAll {
		return reach.NewTrafficContentForAllTraffic()
	}

	// TODO: once new sets are added, finish logic for extracting traffic content

	return reach.TrafficContent{ // and then remove this
		IPProtocol: 0,
		PortSet:    nil,
		ICMPSet:    nil,
	}
}
