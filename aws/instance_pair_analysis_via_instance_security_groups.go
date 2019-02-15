package aws

import (
	"log"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/luhring/cnct/network"
)

type analysisOfInstanceSecurityGroups struct {
	doesAnalysisSuggestThatAccessExists             bool
	isFailureAttributableToLackOfSecurityGroups     bool
	isFailureAttributableToLackOfCommonAllowedPorts bool
	securityGroupIDsOfFirstInstance                 []*string
	securityGroupIDsOfSecondInstance                []*string
	accessiblePorts                                 []*network.PortRange
}

func (instancePair *InstancePair) analyzeNetworkAccessViaInstanceSecurityGroups(ec2Client *ec2.EC2) *analysisOfInstanceSecurityGroups {
	firstInstance := instancePair[0]
	secondInstance := instancePair[1]
	securityGroupsOfFirstInstance, err := getSecurityGroupsForIds(ec2Client, instancePair[0].SecurityGroupIDs)

	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	securityGroupsOfSecondInstance, err := getSecurityGroupsForIds(ec2Client, instancePair[1].SecurityGroupIDs)

	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	if len(securityGroupsOfFirstInstance) == 0 || len(securityGroupsOfSecondInstance) == 0 {
		return &analysisOfInstanceSecurityGroups{
			doesAnalysisSuggestThatAccessExists:             false,
			isFailureAttributableToLackOfSecurityGroups:     true,
			isFailureAttributableToLackOfCommonAllowedPorts: true,
			securityGroupIDsOfFirstInstance:                 instancePair[0].SecurityGroupIDs,
			securityGroupIDsOfSecondInstance:                instancePair[1].SecurityGroupIDs,
			accessiblePorts:                                 nil,
		}
	}

	portsAllowedByFirstInstanceOutboundToSecondInstance :=
		getPortsAllowedOutboundAcrossAllSecurityGroupsForInstanceAccess(
			securityGroupsOfFirstInstance,
			secondInstance,
		)

	portsAllowedBySecondInstanceInboundFromFirstInstance :=
		getPortsAllowedInboundAcrossAllSecurityGroupsForInstanceAccess(
			securityGroupsOfSecondInstance,
			firstInstance,
		)

	portsAllowedByBothInstances := network.GetIntersectionBetweenTwoListsOfPortRanges(
		portsAllowedByFirstInstanceOutboundToSecondInstance,
		portsAllowedBySecondInstanceInboundFromFirstInstance,
	)

	if len(portsAllowedByBothInstances) == 0 {
		return &analysisOfInstanceSecurityGroups{
			doesAnalysisSuggestThatAccessExists:             false,
			isFailureAttributableToLackOfSecurityGroups:     false,
			isFailureAttributableToLackOfCommonAllowedPorts: true,
			securityGroupIDsOfFirstInstance:                 instancePair[0].SecurityGroupIDs,
			securityGroupIDsOfSecondInstance:                instancePair[1].SecurityGroupIDs,
		}
	}

	return &analysisOfInstanceSecurityGroups{
		doesAnalysisSuggestThatAccessExists:             true,
		isFailureAttributableToLackOfSecurityGroups:     false,
		isFailureAttributableToLackOfCommonAllowedPorts: false,
		securityGroupIDsOfFirstInstance:                 instancePair[0].SecurityGroupIDs,
		securityGroupIDsOfSecondInstance:                instancePair[1].SecurityGroupIDs,
		accessiblePorts:                                 portsAllowedByBothInstances,
	}
}

func getPortsAllowedOutboundAcrossAllSecurityGroupsForInstanceAccess(securityGroups []*ec2.SecurityGroup, instance *Instance) []*network.PortRange {
	var allowedPorts []*network.PortRange

	for _, securityGroup := range securityGroups {
		allowedPorts = append(allowedPorts, getPortsAllowedAcrossAllSecurityGroupPermissionsForInstanceAccess(securityGroup.IpPermissionsEgress, instance)...)
	}

	return network.DefragmentPortRanges(allowedPorts)
}

func getPortsAllowedInboundAcrossAllSecurityGroupsForInstanceAccess(securityGroups []*ec2.SecurityGroup, instance *Instance) []*network.PortRange {
	var allowedPorts []*network.PortRange

	for _, securityGroup := range securityGroups {
		allowedPorts = append(allowedPorts, getPortsAllowedAcrossAllSecurityGroupPermissionsForInstanceAccess(securityGroup.IpPermissions, instance)...)
	}

	return network.DefragmentPortRanges(allowedPorts)
}

func getPortsAllowedAcrossAllSecurityGroupPermissionsForInstanceAccess(permissions []*ec2.IpPermission, instance *Instance) []*network.PortRange {
	var allowedPorts []*network.PortRange

	for _, permission := range permissions {
		allowedPorts = append(allowedPorts, getPortsAllowedBySecurityGroupPermissionParametersForInstanceAccess(permission, instance)...)
	}

	return network.DefragmentPortRanges(allowedPorts)
}

func getPortsAllowedBySecurityGroupPermissionParametersForInstanceAccess(permission *ec2.IpPermission, instance *Instance) []*network.PortRange {
	var allowedPorts []*network.PortRange

	for _, ipRange := range permission.IpRanges {
		_, permittedNetwork, _ := net.ParseCIDR(*ipRange.CidrIp)

		for _, privateIPv4Address := range instance.PrivateIPv4Addresses {
			if permittedNetwork.Contains(privateIPv4Address) {
				allowedPorts = append(allowedPorts, getPortRangeFromPermission(permission))
			}
		}

		for _, publicIPv4Address := range instance.PublicIPv4Addresses {
			if permittedNetwork.Contains(publicIPv4Address) {
				allowedPorts = append(allowedPorts, getPortRangeFromPermission(permission))
			}
		}

		for _, ipv6Address := range instance.IPv6Addresses {
			if permittedNetwork.Contains(ipv6Address) {
				allowedPorts = append(allowedPorts, getPortRangeFromPermission(permission))
			}
		}
	}

	if doesSetOfSecurityGroupsForInstanceIntersectWithSecurityGroupsForPermission(instance, permission) {
		allowedPorts = append(allowedPorts, getPortRangeFromPermission(permission))
	}

	return network.DefragmentPortRanges(allowedPorts)
}

func getPortRangeFromPermission(permission *ec2.IpPermission) *network.PortRange {
	doesSpecifyAllProtocols := false
	protocol := *permission.IpProtocol

	if protocol == "-1" {
		doesSpecifyAllProtocols = true
		protocol = ""
	}

	doesSpecifyAllPorts := false
	var lowPort int64
	var highPort int64

	if doesProtocolImplyAllPortsAreAccessible(protocol) {
		doesSpecifyAllPorts = true
	}

	if permission.FromPort == nil && permission.ToPort == nil {
		doesSpecifyAllPorts = true
	} else {
		lowPort = *permission.FromPort
		highPort = *permission.ToPort
	}

	return &network.PortRange{
		DoesSpecifyAllPorts:     doesSpecifyAllPorts,
		LowPort:                 lowPort,
		HighPort:                highPort,
		DoesSpecifyAllProtocols: doesSpecifyAllProtocols,
		Protocol:                protocol,
	}
}

func doesProtocolImplyAllPortsAreAccessible(protocol string) bool {
	return (false == strings.EqualFold(protocol, "tcp") &&
		false == strings.EqualFold(protocol, "udp") &&
		false == strings.EqualFold(protocol, "icmp") &&
		false == strings.EqualFold(protocol, "icmpv6") &&
		protocol != "58")
}

func doesSetOfSecurityGroupsForInstanceIntersectWithSecurityGroupsForPermission(instance *Instance, permission *ec2.IpPermission) bool {
	for _, securityGroupIDForInstance := range instance.SecurityGroupIDs {
		for _, userIDGroupPair := range permission.UserIdGroupPairs {
			if userIDGroupPair.GroupId == securityGroupIDForInstance {
				return true
			}
		}
	}

	return false
}

func (analysis *analysisOfInstanceSecurityGroups) generateExplanationOfLackOfAccess() *string {
	const generalExplanation = "Instances need to have security groups associated with them such that there is at least one port commonly allowed by both the first instance's outbound rules and the second instance's inbound rules."

	explanation := generalExplanation

	if analysis.isFailureAttributableToLackOfSecurityGroups {
		const moreDetailedExplanation = "At least one of the instances has no associated security groups."

		explanation += "\n\n" + moreDetailedExplanation
	}

	return &explanation
}
