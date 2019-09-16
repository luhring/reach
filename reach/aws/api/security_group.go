package api

import (
	"net"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
	reachAWS "github.com/luhring/reach/reach/aws"
)

func (getter *ResourceGetter) GetSecurityGroup(id string) (*reachAWS.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(id),
		},
	}
	result, err := getter.ec2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(result.SecurityGroups, "security group", id); err != nil {
		return nil, err
	}

	securityGroup := newSecurityGroupFromAPI(result.SecurityGroups[0])
	return &securityGroup, nil
}

func newSecurityGroupFromAPI(securityGroup *ec2.SecurityGroup) reachAWS.SecurityGroup {
	inboundRules := getSecurityGroupRules(securityGroup.IpPermissions)
	outboundRules := getSecurityGroupRules(securityGroup.IpPermissionsEgress)

	return reachAWS.SecurityGroup{
		ID:            aws.StringValue(securityGroup.GroupId),
		NameTag:       getNameTag(securityGroup.Tags),
		GroupName:     aws.StringValue(securityGroup.GroupName),
		VPCID:         aws.StringValue(securityGroup.VpcId),
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}
}

func getSecurityGroupRules(inputRules []*ec2.IpPermission) []reachAWS.SecurityGroupRule {
	if inputRules == nil {
		return nil
	}

	rules := make([]reachAWS.SecurityGroupRule, len(inputRules))

	for i, inputRule := range inputRules {
		if inputRule != nil {
			rules[i] = getSecurityGroupRule(inputRule)
		}
	}

	return rules
}

func getSecurityGroupRule(rule *ec2.IpPermission) reachAWS.SecurityGroupRule { // note: this function ignores rule direction (inbound vs. outbound)
	if rule == nil {
		return reachAWS.SecurityGroupRule{}
	}

	trafficContent := getTrafficContentFromSecurityGroupRule(rule)

	// TODO: see if we really need to handle multiple pairs -- the docs don't mention this capability -- https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html#SecurityGroupRules

	var targetSecurityGroupReferenceID, targetSecurityGroupReferenceAccountID string

	if rule.UserIdGroupPairs != nil {
		firstPair := rule.UserIdGroupPairs[0] // if panicking, see above to-do...
		targetSecurityGroupReferenceID = getSecurityGroupReferenceID(firstPair)
		targetSecurityGroupReferenceAccountID = getSecurityGroupReferenceAccountID(firstPair)
	}

	// TODO: Handle prefix lists (and thus VPC endpoints)
	// for context: https://docs.aws.amazon.com/vpc/latest/userguide/vpce-gateway.html

	targetIPNetworks := getIPNetworksFromSecurityGroupRule(rule.IpRanges, rule.Ipv6Ranges)

	return reachAWS.SecurityGroupRule{
		TrafficContent:                        trafficContent,
		TargetSecurityGroupReferenceID:        targetSecurityGroupReferenceID,
		TargetSecurityGroupReferenceAccountID: targetSecurityGroupReferenceAccountID,
		TargetIPNetworks:                      targetIPNetworks,
	}
}

func getTrafficContentFromSecurityGroupRule(rule *ec2.IpPermission) reach.TrafficContent {
	ipProtocolString := aws.StringValue(rule.IpProtocol)
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

func getSecurityGroupReferenceID(pair *ec2.UserIdGroupPair) string {
	if pair == nil {
		return ""
	}

	return aws.StringValue(pair.GroupId)
}

func getSecurityGroupReferenceAccountID(pair *ec2.UserIdGroupPair) string {
	if pair == nil {
		return ""
	}

	return aws.StringValue(pair.UserId)
}

func getIPNetworksFromSecurityGroupRule(ipv4Ranges []*ec2.IpRange, ipv6Ranges []*ec2.Ipv6Range) []*net.IPNet {
	networks := make([]*net.IPNet, len(ipv4Ranges)+len(ipv6Ranges))

	for i, block := range ipv4Ranges {
		if block != nil {
			_, network, err := net.ParseCIDR(aws.StringValue(block.CidrIp))
			if err != nil {
				networks[i] = network
			}
		}
	}

	for i, block := range ipv6Ranges {
		if block != nil {
			_, network, err := net.ParseCIDR(aws.StringValue(block.CidrIpv6))
			if err != nil {
				networks[len(ipv4Ranges)+i] = network
			}
		}
	}

	return networks
}
