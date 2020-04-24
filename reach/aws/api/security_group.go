package api

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/luhring/reach/reach"
	reachAWS "github.com/luhring/reach/reach/aws"
	"github.com/luhring/reach/reach/set"
)

// SecurityGroup queries the AWS API for a security group matching the given ID.
func (client *DomainClient) SecurityGroup(id string) (*reachAWS.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(id),
		},
	}
	result, err := client.ec2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	if err = ensureSingleResult(len(result.SecurityGroups), "security group", id); err != nil {
		return nil, err
	}

	securityGroup := newSecurityGroupFromAPI(result.SecurityGroups[0])
	return &securityGroup, nil
}

func newSecurityGroupFromAPI(securityGroup *ec2.SecurityGroup) reachAWS.SecurityGroup {
	inboundRules := securityGroupRules(securityGroup.IpPermissions)
	outboundRules := securityGroupRules(securityGroup.IpPermissionsEgress)

	return reachAWS.SecurityGroup{
		ID:            aws.StringValue(securityGroup.GroupId),
		NameTag:       nameTag(securityGroup.Tags),
		GroupName:     aws.StringValue(securityGroup.GroupName),
		VPCID:         aws.StringValue(securityGroup.VpcId),
		InboundRules:  inboundRules,
		OutboundRules: outboundRules,
	}
}

func securityGroupRules(inputRules []*ec2.IpPermission) []reachAWS.SecurityGroupRule {
	if inputRules == nil {
		return nil
	}

	rules := make([]reachAWS.SecurityGroupRule, len(inputRules))

	for i, inputRule := range inputRules {
		if inputRule != nil {
			rules[i] = securityGroupRule(inputRule)
		}
	}

	return rules
}

func securityGroupRule(rule *ec2.IpPermission) reachAWS.SecurityGroupRule { // note: this function ignores rule direction (inbound vs. outbound)
	if rule == nil {
		return reachAWS.SecurityGroupRule{}
	}

	tc, err := trafficContentFromAWSIPPermission(rule)
	if err != nil {
		panic(err) // TODO: Better error handling
	}

	// TODO: see if we really need to handle multiple pairs -- the docs don't mention this capability -- https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html#SecurityGroupRules

	var targetSecurityGroupReferenceID, targetSecurityGroupReferenceAccountID string

	if rule.UserIdGroupPairs != nil {
		firstPair := rule.UserIdGroupPairs[0] // if panicking, see above to-do...
		targetSecurityGroupReferenceID = securityGroupReferenceID(firstPair)
		targetSecurityGroupReferenceAccountID = securityGroupReferenceAccountID(firstPair)
	}

	// TODO: Handle prefix lists (and thus VPC endpoints)
	// for context: https://docs.aws.amazon.com/vpc/latest/userguide/vpce-gateway.html

	targetIPNetworks := ipNetworksFromSecurityGroupRule(rule.IpRanges, rule.Ipv6Ranges)

	return reachAWS.SecurityGroupRule{
		TrafficContent:                        tc,
		TargetSecurityGroupReferenceID:        targetSecurityGroupReferenceID,
		TargetSecurityGroupReferenceAccountID: targetSecurityGroupReferenceAccountID,
		TargetIPNetworks:                      targetIPNetworks,
	}
}

func newPortSetFromAWSPortRange(portRange *ec2.PortRange) (set.PortSet, error) {
	if portRange == nil {
		return set.PortSet{}, fmt.Errorf("input portRange was nil")
	}

	from := aws.Int64Value(portRange.From)
	to := aws.Int64Value(portRange.To)

	return set.NewPortSetFromRange(uint16(from), uint16(to))
}

func newPortSetFromAWSIPPermission(permission *ec2.IpPermission) (set.PortSet, error) {
	if permission == nil {
		return set.PortSet{}, fmt.Errorf("input IpPermission was nil")
	}

	from := aws.Int64Value(permission.FromPort)
	to := aws.Int64Value(permission.ToPort)

	return set.NewPortSetFromRange(uint16(from), uint16(to))
}

func securityGroupReferenceID(pair *ec2.UserIdGroupPair) string {
	if pair == nil {
		return ""
	}

	return aws.StringValue(pair.GroupId)
}

func securityGroupReferenceAccountID(pair *ec2.UserIdGroupPair) string {
	if pair == nil {
		return ""
	}

	return aws.StringValue(pair.UserId)
}

func ipNetworksFromSecurityGroupRule(ipv4Ranges []*ec2.IpRange, ipv6Ranges []*ec2.Ipv6Range) []net.IPNet {
	networks := make([]net.IPNet, len(ipv4Ranges)+len(ipv6Ranges))

	for i, block := range ipv4Ranges {
		if block != nil {
			_, network, err := net.ParseCIDR(aws.StringValue(block.CidrIp))
			if err == nil {
				networks[i] = *network
			}
		}
	}

	for i, block := range ipv6Ranges {
		if block != nil {
			_, network, err := net.ParseCIDR(aws.StringValue(block.CidrIpv6))
			if err == nil {
				networks[len(ipv4Ranges)+i] = *network
			}
		}
	}

	return networks
}

func trafficContentFromAWSIPPermission(permission *ec2.IpPermission) (reach.TrafficContent, error) {
	const errCreation = "unable to create content: %v"

	protocol, err := convertAWSIPProtocolStringToProtocol(permission.IpProtocol)
	if err != nil {
		return reach.TrafficContent{}, fmt.Errorf(errCreation, err)
	}

	if protocol == reach.ProtocolAll {
		return reach.NewTrafficContentForAllTraffic(), nil
	}

	if protocol.UsesPorts() {
		portSet, err := newPortSetFromAWSIPPermission(permission)
		if err != nil {
			return reach.TrafficContent{}, fmt.Errorf(errCreation, err)
		}

		return reach.NewTrafficContentForPorts(protocol, portSet), nil
	}

	if protocol == reach.ProtocolICMPv4 || protocol == reach.ProtocolICMPv6 {
		icmpSet, err := newICMPSetFromAWSIPPermission(permission)
		if err != nil {
			return reach.TrafficContent{}, fmt.Errorf(errCreation, err)
		}

		return reach.NewTrafficContentForICMP(protocol, icmpSet), nil
	}

	return reach.NewTrafficContentForCustomProtocol(protocol, true), nil
}

func newICMPSetFromAWSICMPTypeCode(icmpTypeCode *ec2.IcmpTypeCode) (set.ICMPSet, error) {
	if icmpTypeCode == nil {
		return set.ICMPSet{}, fmt.Errorf("input icmpTypeCode was nil")
	}

	icmpType := aws.Int64Value(icmpTypeCode.Type)

	if icmpType == set.AllICMPTypes {
		result := set.NewFullICMPSet()
		return result, nil
	}

	icmpTypeValue := uint8(icmpType) // i.e. equivalent to ICMP header value

	icmpCode := aws.Int64Value(icmpTypeCode.Code)

	if icmpCode == set.AllICMPCodes {
		return set.NewICMPSetFromICMPType(icmpTypeValue)
	}

	icmpCodeValue := uint8(icmpCode) // i.e. equivalent to ICMP header value

	return set.NewICMPSetFromICMPTypeCode(icmpTypeValue, icmpCodeValue)
}

func newICMPSetFromAWSIPPermission(permission *ec2.IpPermission) (set.ICMPSet, error) {
	if permission == nil {
		return set.ICMPSet{}, fmt.Errorf("input IpPermission was nil")
	}

	icmpType := aws.Int64Value(permission.FromPort)

	if icmpType == set.AllICMPTypes {
		result := set.NewFullICMPSet()
		return result, nil
	}

	icmpTypeValue := uint8(icmpType) // i.e. equivalent to ICMP header value

	icmpCode := aws.Int64Value(permission.ToPort)

	if icmpCode == set.AllICMPCodes {
		return set.NewICMPSetFromICMPType(icmpTypeValue)
	}

	icmpCodeValue := uint8(icmpCode) // i.e. equivalent to ICMP header value

	return set.NewICMPSetFromICMPTypeCode(icmpTypeValue, icmpCodeValue)
}
