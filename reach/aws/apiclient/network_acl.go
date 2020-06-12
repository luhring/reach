package apiclient

import (
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	reachAWS "github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/reacherr"
	"github.com/luhring/reach/reach/traffic"
)

// NetworkACL queries the AWS API for a network ACL matching the given ID.
func (client *DomainClient) NetworkACL(id string) (*reachAWS.NetworkACL, error) {
	if r := client.cachedResource(reachAWS.NetworkACLRef(id)); r != nil {
		if v, ok := r.(*reachAWS.NetworkACL); ok {
			return v, nil
		}
	}

	input := &ec2.DescribeNetworkAclsInput{
		NetworkAclIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeNetworkAcls(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, reacherr.New(err, awsErrMessage(aerr))
		}
		return nil, err
	}

	if err = ensureSingleResult(len(result.NetworkAcls), "network ACL", id); err != nil {
		return nil, err
	}

	networkACL, err := newNetworkACLFromAPI(result.NetworkAcls[0])
	if err != nil {
		return nil, err
	}
	client.cacheResource(networkACL)
	return &networkACL, nil
}

func newNetworkACLFromAPI(networkACL *ec2.NetworkAcl) (reachAWS.NetworkACL, error) {
	inboundRules, err := networkACLRulesForDirection(networkACL.Entries, reachAWS.NetworkACLRuleDirectionInbound)
	if err != nil {
		return reachAWS.NetworkACL{}, err
	}
	outboundRules, err := networkACLRulesForDirection(networkACL.Entries, reachAWS.NetworkACLRuleDirectionOutbound)
	if err != nil {
		return reachAWS.NetworkACL{}, err
	}

	return reachAWS.NetworkACL{
		ID:            aws.StringValue(networkACL.NetworkAclId),
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}, nil
}

func networkACLRulesForDirection(entries []*ec2.NetworkAclEntry, direction reachAWS.NetworkACLRuleDirection) ([]reachAWS.NetworkACLRule, error) {
	if entries == nil {
		return nil, nil
	}

	var rules []reachAWS.NetworkACLRule

	for _, entry := range entries {
		if entry != nil {
			if directionMatches(direction, *entry) {
				rule, err := networkACLRule(*entry)
				if err != nil {
					return nil, err
				}
				rules = append(rules, rule)
			}
		}
	}

	return rules, nil
}

func directionMatches(direction reachAWS.NetworkACLRuleDirection, entry ec2.NetworkAclEntry) bool {
	outboundEntry := aws.BoolValue(entry.Egress)
	return outboundEntry == (direction == reachAWS.NetworkACLRuleDirectionOutbound)
}

func networkACLRule(entry ec2.NetworkAclEntry) (reachAWS.NetworkACLRule, error) { // note: this function ignores rule direction (inbound vs. outbound)
	_, targetIPNetwork, err := net.ParseCIDR(aws.StringValue(entry.CidrBlock))
	if err != nil {
		return reachAWS.NetworkACLRule{}, err
	}

	var action reachAWS.NetworkACLRuleAction

	if aws.StringValue(entry.RuleAction) == ec2.RuleActionAllow {
		action = reachAWS.NetworkACLRuleActionAllow
	} else {
		action = reachAWS.NetworkACLRuleActionDeny
	}

	tc, err := newTrafficContentFromAWSNACLEntry(entry)
	if err != nil {
		return reachAWS.NetworkACLRule{}, err
	}

	return reachAWS.NetworkACLRule{
		Number:          aws.Int64Value(entry.RuleNumber),
		TrafficContent:  tc,
		TargetIPNetwork: targetIPNetwork,
		Action:          action,
	}, nil
}

func newTrafficContentFromAWSNACLEntry(entry ec2.NetworkAclEntry) (traffic.Content, error) {
	protocol, err := convertAWSIPProtocolStringToProtocol(entry.Protocol)
	if err != nil {
		return traffic.Content{}, err
	}

	if protocol == traffic.ProtocolAll {
		return traffic.All(), nil
	}

	if protocol.UsesPorts() {
		portSet, err := newPortSetFromAWSPortRange(entry.PortRange)
		if err != nil {
			return traffic.Content{}, err
		}

		return traffic.ForPorts(protocol, portSet), nil
	}

	if protocol == traffic.ProtocolICMPv4 || protocol == traffic.ProtocolICMPv6 {
		icmpSet, err := newICMPSetFromAWSICMPTypeCode(entry.IcmpTypeCode)
		if err != nil {
			return traffic.Content{}, err
		}

		return traffic.ForICMP(protocol, icmpSet), nil
	}

	return traffic.ForCustomProtocol(protocol, true), nil
}
